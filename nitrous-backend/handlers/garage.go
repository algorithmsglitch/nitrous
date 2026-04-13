package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ── NHTSA response types ──────────────────────────────────────────────────────
// NHTSA vPIC API: https://vpic.nhtsa.dot.gov/api/

type NHTSAMake struct {
	MakeID   int    `json:"Make_ID"`
	MakeName string `json:"Make_Name"`
}

type NHTSAMakesResponse struct {
	Results []NHTSAMake `json:"Results"`
}

type NHTSAModel struct {
	ModelID   int    `json:"Model_ID"`
	ModelName string `json:"Model_Name"`
	MakeID    int    `json:"Make_ID"`
	MakeName  string `json:"Make_Name"`
}

type NHTSAModelsResponse struct {
	Results []NHTSAModel `json:"Results"`
}

// ── Domain types ──────────────────────────────────────────────────────────────

type VehicleSpec struct {
	Make         string  `json:"make"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Trim         string  `json:"trim"`
	Engine       string  `json:"engine"`
	Displacement int     `json:"displacement"`
	Cylinders    int     `json:"cylinders"`
	HP           float64 `json:"hp"`
	Torque       float64 `json:"torque"`
	TopSpeed     float64 `json:"topSpeed"`
	Weight       float64 `json:"weight"`
	ZeroToSixty  float64 `json:"zeroToSixty"`
	Drivetrain   string  `json:"drivetrain"`
	FuelType     string  `json:"fuelType"`
	Seats        int     `json:"seats"`
}

type TunedStats struct {
	HP          float64 `json:"hp"`
	Torque      float64 `json:"torque"`
	TopSpeed    float64 `json:"topSpeed"`
	ZeroToSixty float64 `json:"zeroToSixty"`
	Weight      float64 `json:"weight"`
	Config      string  `json:"config"`
}

type Delta struct {
	HP          float64 `json:"hp"`
	Torque      float64 `json:"torque"`
	TopSpeed    float64 `json:"topSpeed"`
	ZeroToSixty float64 `json:"zeroToSixty"`
	Weight      float64 `json:"weight"`
}

type TuneResponse struct {
	Base   VehicleSpec  `json:"base"`
	Tuned  TunedStats   `json:"tuned"`
	Delta  Delta        `json:"delta"`
	Config TuningConfig `json:"config"`
}

type TuneRequest struct {
	Make   string `json:"make"   binding:"required"`
	Model  string `json:"model"  binding:"required"`
	Year   int    `json:"year"   binding:"required"`
	Tuning string `json:"tuning" binding:"required"`
}

// ── Tuning configs ────────────────────────────────────────────────────────────

type TuningConfig struct {
	Label        string  `json:"label"`
	HPMult       float64 `json:"hpMult"`
	TorqueMult   float64 `json:"torqueMult"`
	TopSpeedMult float64 `json:"topSpeedMult"`
	ZeroMult     float64 `json:"zeroMult"`
	WeightMult   float64 `json:"weightMult"`
}

var tuningConfigs = map[string]TuningConfig{
	"stock":  {Label: "Stock", HPMult: 1.00, TorqueMult: 1.00, TopSpeedMult: 1.00, ZeroMult: 1.00, WeightMult: 1.00},
	"street": {Label: "Street", HPMult: 1.08, TorqueMult: 1.06, TopSpeedMult: 1.04, ZeroMult: 0.95, WeightMult: 0.97},
	"track":  {Label: "Track", HPMult: 1.18, TorqueMult: 1.12, TopSpeedMult: 1.10, ZeroMult: 0.86, WeightMult: 0.90},
	"race":   {Label: "Race Spec", HPMult: 1.35, TorqueMult: 1.25, TopSpeedMult: 1.18, ZeroMult: 0.76, WeightMult: 0.82},
	"drift":  {Label: "Drift", HPMult: 1.20, TorqueMult: 1.30, TopSpeedMult: 0.96, ZeroMult: 0.92, WeightMult: 0.94},
}

// ── Curated vehicle database ──────────────────────────────────────────────────
// Since CarQuery is unreliable, all performance data lives here.
// Sources: manufacturer press kits, verified road tests, Motor Trend, Car and Driver.

type vehicleKey struct {
	Make  string
	Model string
	Year  int
}

type vehicleData struct {
	Trim         string
	Engine       string
	Displacement int
	Cylinders    int
	HP           float64
	Torque       float64 // lb-ft
	TopSpeed     float64 // mph
	Weight       float64 // lbs
	ZeroToSixty  float64 // seconds
	Drivetrain   string
	FuelType     string
	Seats        int
}

var vehicleDB = map[vehicleKey]vehicleData{
	// ── Ferrari ──
	{Make: "ferrari", Model: "f40", Year: 1992}: {
		Trim: "Base", Engine: "Twin-Turbo V8", Displacement: 2936, Cylinders: 8,
		HP: 478, Torque: 426, TopSpeed: 201, Weight: 2425, ZeroToSixty: 3.8,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	{Make: "ferrari", Model: "f8 tributo", Year: 2022}: {
		Trim: "Base", Engine: "Twin-Turbo V8", Displacement: 3902, Cylinders: 8,
		HP: 710, Torque: 568, TopSpeed: 211, Weight: 3164, ZeroToSixty: 2.9,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	// ── Porsche ──
	{Make: "porsche", Model: "911 gt3 rs", Year: 2024}: {
		Trim: "GT3 RS", Engine: "Naturally Aspirated Flat-6", Displacement: 3996, Cylinders: 6,
		HP: 518, Torque: 343, TopSpeed: 184, Weight: 3268, ZeroToSixty: 3.0,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	{Make: "porsche", Model: "911 turbo s", Year: 2024}: {
		Trim: "Turbo S", Engine: "Twin-Turbo Flat-6", Displacement: 3745, Cylinders: 6,
		HP: 640, Torque: 590, TopSpeed: 205, Weight: 3627, ZeroToSixty: 2.6,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── Nissan ──
	{Make: "nissan", Model: "gt-r", Year: 2023}: {
		Trim: "Premium", Engine: "Twin-Turbo V6", Displacement: 3799, Cylinders: 6,
		HP: 565, Torque: 467, TopSpeed: 196, Weight: 3927, ZeroToSixty: 2.9,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 4,
	},
	{Make: "nissan", Model: "z", Year: 2024}: {
		Trim: "Performance", Engine: "Twin-Turbo V6", Displacement: 3000, Cylinders: 6,
		HP: 400, Torque: 350, TopSpeed: 155, Weight: 3306, ZeroToSixty: 4.5,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	// ── Lamborghini ──
	{Make: "lamborghini", Model: "huracan evo", Year: 2023}: {
		Trim: "EVO", Engine: "Naturally Aspirated V10", Displacement: 5204, Cylinders: 10,
		HP: 630, Torque: 443, TopSpeed: 202, Weight: 3135, ZeroToSixty: 2.9,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 2,
	},
	{Make: "lamborghini", Model: "urus", Year: 2023}: {
		Trim: "S", Engine: "Twin-Turbo V8", Displacement: 3996, Cylinders: 8,
		HP: 657, Torque: 627, TopSpeed: 189, Weight: 4850, ZeroToSixty: 3.5,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 5,
	},
	// ── Toyota ──
	{Make: "toyota", Model: "gr supra", Year: 2024}: {
		Trim: "3.0 A91-MT", Engine: "Turbocharged I6", Displacement: 2998, Cylinders: 6,
		HP: 382, Torque: 368, TopSpeed: 155, Weight: 3181, ZeroToSixty: 3.9,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	{Make: "toyota", Model: "gr86", Year: 2024}: {
		Trim: "Premium", Engine: "Naturally Aspirated Flat-4", Displacement: 2387, Cylinders: 4,
		HP: 228, Torque: 184, TopSpeed: 140, Weight: 2822, ZeroToSixty: 6.1,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── Lotus ──
	{Make: "lotus", Model: "evora gt", Year: 2022}: {
		Trim: "GT", Engine: "Supercharged V6", Displacement: 3456, Cylinders: 6,
		HP: 416, Torque: 317, TopSpeed: 188, Weight: 3046, ZeroToSixty: 3.8,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── BMW ──
	{Make: "bmw", Model: "m3 competition", Year: 2024}: {
		Trim: "Competition", Engine: "Twin-Turbo I6", Displacement: 2993, Cylinders: 6,
		HP: 503, Torque: 479, TopSpeed: 180, Weight: 3868, ZeroToSixty: 3.4,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 5,
	},
	{Make: "bmw", Model: "m4 competition", Year: 2024}: {
		Trim: "Competition", Engine: "Twin-Turbo I6", Displacement: 2993, Cylinders: 6,
		HP: 503, Torque: 479, TopSpeed: 180, Weight: 3814, ZeroToSixty: 3.4,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── Ford ──
	{Make: "ford", Model: "mustang gt500", Year: 2023}: {
		Trim: "Shelby GT500", Engine: "Supercharged V8", Displacement: 5163, Cylinders: 8,
		HP: 760, Torque: 625, TopSpeed: 180, Weight: 4225, ZeroToSixty: 3.3,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	{Make: "ford", Model: "mustang dark horse", Year: 2024}: {
		Trim: "Dark Horse", Engine: "Naturally Aspirated V8", Displacement: 5038, Cylinders: 8,
		HP: 500, Torque: 418, TopSpeed: 155, Weight: 4060, ZeroToSixty: 4.0,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── Chevrolet ──
	{Make: "chevrolet", Model: "corvette z06", Year: 2024}: {
		Trim: "Z06", Engine: "Naturally Aspirated Flat-Plane V8", Displacement: 5497, Cylinders: 8,
		HP: 670, Torque: 460, TopSpeed: 196, Weight: 3366, ZeroToSixty: 2.6,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	{Make: "chevrolet", Model: "camaro zl1", Year: 2024}: {
		Trim: "ZL1", Engine: "Supercharged V8", Displacement: 6162, Cylinders: 8,
		HP: 650, Torque: 650, TopSpeed: 185, Weight: 4120, ZeroToSixty: 3.5,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 4,
	},
	// ── McLaren ──
	{Make: "mclaren", Model: "720s", Year: 2023}: {
		Trim: "Base", Engine: "Twin-Turbo V8", Displacement: 3994, Cylinders: 8,
		HP: 710, Torque: 568, TopSpeed: 212, Weight: 2937, ZeroToSixty: 2.8,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
	// ── Dodge ──
	{Make: "dodge", Model: "challenger hellcat", Year: 2023}: {
		Trim: "SRT Hellcat", Engine: "Supercharged V8", Displacement: 6166, Cylinders: 8,
		HP: 717, Torque: 656, TopSpeed: 199, Weight: 4439, ZeroToSixty: 3.6,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 5,
	},
	// ── Subaru ──
	{Make: "subaru", Model: "wrx sti", Year: 2021}: {
		Trim: "STI", Engine: "Turbocharged Boxer-4", Displacement: 2457, Cylinders: 4,
		HP: 310, Torque: 290, TopSpeed: 159, Weight: 3395, ZeroToSixty: 4.7,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 5,
	},
	// ── Audi ──
	{Make: "audi", Model: "r8 v10", Year: 2023}: {
		Trim: "V10 Performance", Engine: "Naturally Aspirated V10", Displacement: 5204, Cylinders: 10,
		HP: 602, Torque: 413, TopSpeed: 205, Weight: 3571, ZeroToSixty: 3.1,
		Drivetrain: "AWD", FuelType: "Gasoline", Seats: 2,
	},
	// ── Mercedes ──
	{Make: "mercedes-benz", Model: "amg gt black series", Year: 2021}: {
		Trim: "AMG GT Black Series", Engine: "Twin-Turbo V8", Displacement: 3982, Cylinders: 8,
		HP: 720, Torque: 590, TopSpeed: 202, Weight: 3638, ZeroToSixty: 3.1,
		Drivetrain: "RWD", FuelType: "Gasoline", Seats: 2,
	},
}

// ── NHTSA client ──────────────────────────────────────────────────────────────

const nhtsaBase = "https://vpic.nhtsa.dot.gov/api/vehicles"

var garageHTTPClient = &http.Client{Timeout: 8 * time.Second}

func nhtsaFetch(endpoint string) ([]byte, error) {
	resp, err := garageHTTPClient.Get(nhtsaBase + endpoint + "?format=json")
	if err != nil {
		return nil, fmt.Errorf("nhtsa request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nhtsa read: %w", err)
	}
	return body, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func normaliseKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func lookupVehicle(make_, model string, year int) (vehicleData, bool) {
	key := vehicleKey{
		Make:  normaliseKey(make_),
		Model: normaliseKey(model),
		Year:  year,
	}
	data, ok := vehicleDB[key]
	return data, ok
}

func applyTuning(base VehicleSpec, cfg TuningConfig) TunedStats {
	return TunedStats{
		HP:          math.Round(base.HP * cfg.HPMult),
		Torque:      math.Round(base.Torque * cfg.TorqueMult),
		TopSpeed:    math.Round(base.TopSpeed * cfg.TopSpeedMult),
		ZeroToSixty: round2(base.ZeroToSixty * cfg.ZeroMult),
		Weight:      math.Round(base.Weight * cfg.WeightMult),
		Config:      cfg.Label,
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GET /api/garage/makes
// Returns all makes from NHTSA vPIC (all vehicle types).
func handleMakes(c *gin.Context) {
	raw, err := nhtsaFetch("/getallmakes")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var resp NHTSAMakesResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	type Make struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	out := make([]Make, 0, len(resp.Results))
	for _, m := range resp.Results {
		out = append(out, Make{ID: m.MakeID, Name: m.MakeName})
	}
	c.JSON(http.StatusOK, gin.H{"makes": out})
}

// GET /api/garage/models?make=Ferrari&year=2024
func handleModels(c *gin.Context) {
	make_ := strings.TrimSpace(c.Query("make"))
	if make_ == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make is required"})
		return
	}
	year := c.Query("year")
	var endpoint string
	if year != "" {
		endpoint = fmt.Sprintf("/getmodelsformakeyear/make/%s/modelyear/%s/vehicleType/car",
			make_, year)
	} else {
		endpoint = fmt.Sprintf("/getmodelsformake/%s", make_)
	}

	raw, err := nhtsaFetch(endpoint)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var resp NHTSAModelsResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	type Model struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Make string `json:"make"`
	}
	out := make([]Model, 0, len(resp.Results))
	for _, m := range resp.Results {
		out = append(out, Model{ID: m.ModelID, Name: m.ModelName, Make: m.MakeName})
	}
	c.JSON(http.StatusOK, gin.H{"models": out})
}

// GET /api/garage/vehicle?make=Ferrari&model=F40&year=1992
func handleVehicle(c *gin.Context) {
	make_ := c.Query("make")
	model := c.Query("model")
	yearS := c.Query("year")
	if make_ == "" || model == "" || yearS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make, model, and year are required"})
		return
	}
	year, err := strconv.Atoi(yearS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	data, ok := lookupVehicle(make_, model, year)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "vehicle not in database",
			"queried_make":  normaliseKey(make_),
			"queried_model": normaliseKey(model),
			"queried_year":  year,
			"hint":          "add this vehicle to vehicleDB in garage_handlers.go",
		})
		return
	}

	spec := VehicleSpec{
		Make: make_, Model: model, Year: year,
		Trim: data.Trim, Engine: data.Engine,
		Displacement: data.Displacement, Cylinders: data.Cylinders,
		HP: data.HP, Torque: data.Torque, TopSpeed: data.TopSpeed,
		Weight: data.Weight, ZeroToSixty: data.ZeroToSixty,
		Drivetrain: data.Drivetrain, FuelType: data.FuelType, Seats: data.Seats,
	}
	c.JSON(http.StatusOK, gin.H{"vehicle": spec})
}

// GET /api/garage/tuning-configs
func handleTuningConfigs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"configs": tuningConfigs})
}

// POST /api/garage/tune
func handleTune(c *gin.Context) {
	var req TuneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, ok := tuningConfigs[req.Tuning]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown tuning config: " + req.Tuning})
		return
	}

	data, ok := lookupVehicle(req.Make, req.Model, req.Year)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not in database"})
		return
	}

	base := VehicleSpec{
		Make: req.Make, Model: req.Model, Year: req.Year,
		Trim: data.Trim, Engine: data.Engine,
		Displacement: data.Displacement, Cylinders: data.Cylinders,
		HP: data.HP, Torque: data.Torque, TopSpeed: data.TopSpeed,
		Weight: data.Weight, ZeroToSixty: data.ZeroToSixty,
		Drivetrain: data.Drivetrain, FuelType: data.FuelType, Seats: data.Seats,
	}
	tuned := applyTuning(base, cfg)
	delta := Delta{
		HP:          tuned.HP - base.HP,
		Torque:      tuned.Torque - base.Torque,
		TopSpeed:    tuned.TopSpeed - base.TopSpeed,
		ZeroToSixty: round2(base.ZeroToSixty - tuned.ZeroToSixty),
		Weight:      tuned.Weight - base.Weight,
	}

	c.JSON(http.StatusOK, TuneResponse{Base: base, Tuned: tuned, Delta: delta, Config: cfg})
}

// GET /api/garage/search?q=ferrari
func handleSearch(c *gin.Context) {
	q := normaliseKey(c.Query("q"))
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}
	type Result struct {
		Make  string `json:"make"`
		Model string `json:"model"`
		Year  int    `json:"year"`
	}
	var results []Result
	for k := range vehicleDB {
		if strings.Contains(k.Make, q) || strings.Contains(k.Model, q) {
			results = append(results, Result{Make: k.Make, Model: k.Model, Year: k.Year})
			if len(results) >= 10 {
				break
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// ── Router setup ──────────────────────────────────────────────────────────────

func RegisterGarageRoutes(r *gin.Engine) {
	garage := r.Group("/api/garage")
	{
		garage.GET("/makes", handleMakes)
		garage.GET("/models", handleModels)
		garage.GET("/vehicle", handleVehicle)
		garage.GET("/tuning-configs", handleTuningConfigs)
		garage.POST("/tune", handleTune)
		garage.GET("/search", handleSearch)
	}
}
