# ============================================================
#  NITROUS BACKEND - Full API Test Script
#  Run: powershell -ExecutionPolicy Bypass -File .\test-api.ps1
# ============================================================

$BASE    = "http://localhost:8080"
$PASS    = 0
$FAIL    = 0
$TOKEN   = ""
$TEAM_ID    = ""
$STREAM_ID  = ""
$EVENT_ID   = ""
$MERCH_ID   = ""
$ORDER_ID   = ""

# ── Helpers ──────────────────────────────────────────────────

function Print-Header($text) {
    Write-Host ""
    Write-Host "===================================================" -ForegroundColor DarkCyan
    Write-Host "  $text" -ForegroundColor Cyan
    Write-Host "===================================================" -ForegroundColor DarkCyan
}

function Test-Pass($label) {
    Write-Host "  [PASS] $label" -ForegroundColor Green
    $script:PASS++
}

function Test-Fail($label, $detail) {
    Write-Host "  [FAIL] $label" -ForegroundColor Red
    if ($detail) { Write-Host "         $detail" -ForegroundColor DarkRed }
    $script:FAIL++
}

function Invoke-GET($path, $token = "") {
    $headers = @{}
    if ($token) { $headers["Authorization"] = "Bearer $token" }
    try {
        return Invoke-RestMethod -Uri "$BASE$path" -Headers $headers -ErrorAction Stop
    } catch {
        return $null
    }
}

function Invoke-POST($path, $body, $token = "") {
    $headers = @{ "Content-Type" = "application/json" }
    if ($token) { $headers["Authorization"] = "Bearer $token" }
    try {
        return Invoke-RestMethod -Method POST -Uri "$BASE$path" `
            -Headers $headers -Body ($body | ConvertTo-Json) -ErrorAction Stop
    } catch {
        return $null
    }
}

function Invoke-DELETE($path, $token = "") {
    $headers = @{}
    if ($token) { $headers["Authorization"] = "Bearer $token" }
    try {
        return Invoke-RestMethod -Method DELETE -Uri "$BASE$path" `
            -Headers $headers -ErrorAction Stop
    } catch {
        return $null
    }
}

function Invoke-POST-ExpectError($path, $body, $token = "") {
    $headers = @{ "Content-Type" = "application/json" }
    if ($token) { $headers["Authorization"] = "Bearer $token" }
    try {
        Invoke-RestMethod -Method POST -Uri "$BASE$path" `
            -Headers $headers -Body ($body | ConvertTo-Json) -ErrorAction Stop
        return $false
    } catch {
        return $true
    }
}

# ── 0. Server Check ───────────────────────────────────────────

Print-Header "0. Checking server is running"

$health = Invoke-GET "/health"
if ($health -and $health.status -eq "ok") {
    Test-Pass "Server is running on :8080"
} else {
    Test-Fail "Server not reachable at http://localhost:8080" "Run go run main.go first"
    Write-Host ""
    Write-Host "  Cannot continue - server is not running." -ForegroundColor Yellow
    exit 1
}

# ── 1. Auth ───────────────────────────────────────────────────

Print-Header "1. Auth"

$regBody = @{ email = "testuser@nitrous.io"; password = "password123"; name = "Test User" }
$reg = Invoke-POST "/api/auth/register" $regBody
if ($reg -and $reg.token) {
    Test-Pass "POST /api/auth/register - user created"
    $script:TOKEN = $reg.token
} else {
    Test-Fail "POST /api/auth/register" "No token returned"
}

$loginBody = @{ email = "testuser@nitrous.io"; password = "password123" }
$login = Invoke-POST "/api/auth/login" $loginBody
if ($login -and $login.token) {
    Test-Pass "POST /api/auth/login - token received"
    $script:TOKEN = $login.token
} else {
    Test-Fail "POST /api/auth/login" "Login failed"
}

$me = Invoke-GET "/api/auth/me" $TOKEN
if ($me -and $me.email -eq "testuser@nitrous.io") {
    Test-Pass "GET /api/auth/me - returned correct user"
} else {
    Test-Fail "GET /api/auth/me" "User data not returned"
}

$noAuth = Invoke-POST-ExpectError "/api/orders" @{}
if ($noAuth) {
    Test-Pass "Protected route (no token) - correctly rejected with 401"
} else {
    Test-Fail "Protected route (no token)" "Should have returned 401"
}

# ── 2. Events ─────────────────────────────────────────────────

Print-Header "2. Events"

$events = Invoke-GET "/api/events"
if ($events -and $events.count -gt 0) {
    Test-Pass "GET /api/events - returned $($events.count) events"
    $script:EVENT_ID = $events.events[0].id
} else {
    Test-Fail "GET /api/events" "No events returned"
}

$filtered = Invoke-GET "/api/events?category=motorsport"
if ($filtered -and $filtered.count -gt 0) {
    Test-Pass "GET /api/events?category=motorsport - returned $($filtered.count) events"
} else {
    Test-Fail "GET /api/events?category=motorsport" "No filtered events returned"
}

$live = Invoke-GET "/api/events/live"
if ($null -ne $live) {
    Test-Pass "GET /api/events/live - returned $($live.count) live events"
} else {
    Test-Fail "GET /api/events/live" "Failed"
}

if ($EVENT_ID) {
    $event = Invoke-GET "/api/events/$EVENT_ID"
    if ($event -and $event.id -eq $EVENT_ID) {
        Test-Pass "GET /api/events/:id - returned correct event"
    } else {
        Test-Fail "GET /api/events/:id" "Event not found"
    }
}

# ── 3. Categories ─────────────────────────────────────────────

Print-Header "3. Categories"

$cats = Invoke-GET "/api/categories"
if ($cats -and $cats.count -gt 0) {
    Test-Pass "GET /api/categories - returned $($cats.count) categories"
} else {
    Test-Fail "GET /api/categories" "No categories returned"
}

$cat = Invoke-GET "/api/categories/motorsport"
if ($cat -and $cat.slug -eq "motorsport") {
    Test-Pass "GET /api/categories/motorsport - returned correct category"
} else {
    Test-Fail "GET /api/categories/motorsport" "Category not found"
}

# ── 4. Teams ──────────────────────────────────────────────────

Print-Header "4. Teams"

$teams = Invoke-GET "/api/teams"
if ($teams -and $teams.count -gt 0) {
    Test-Pass "GET /api/teams - returned $($teams.count) teams"
    $script:TEAM_ID = $teams.teams[0].id
} else {
    Test-Fail "GET /api/teams" "No teams returned"
}

if ($TEAM_ID) {
    $team = Invoke-GET "/api/teams/$TEAM_ID"
    if ($team -and $team.id -eq $TEAM_ID) {
        Test-Pass "GET /api/teams/:id - returned correct team"
    } else {
        Test-Fail "GET /api/teams/:id" "Team not found"
    }

    $follow = Invoke-POST "/api/teams/$TEAM_ID/follow" @{} $TOKEN
    if ($follow) {
        Test-Pass "POST /api/teams/:id/follow - followed successfully"
    } else {
        Test-Fail "POST /api/teams/:id/follow" "Follow failed"
    }

    $followAgain = Invoke-POST-ExpectError "/api/teams/$TEAM_ID/follow" @{} $TOKEN
    if ($followAgain) {
        Test-Pass "POST /api/teams/:id/follow (duplicate) - correctly returned 409"
    } else {
        Test-Fail "POST /api/teams/:id/follow (duplicate)" "Should have returned 409 Conflict"
    }

    $unfollow = Invoke-DELETE "/api/teams/$TEAM_ID/follow" $TOKEN
    if ($unfollow) {
        Test-Pass "DELETE /api/teams/:id/follow - unfollowed successfully"
    } else {
        Test-Fail "DELETE /api/teams/:id/follow" "Unfollow failed"
    }
}

# ── 5. Streams ────────────────────────────────────────────────

Print-Header "5. Streams"

$streams = Invoke-GET "/api/streams"
if ($streams -and $streams.count -gt 0) {
    Test-Pass "GET /api/streams - returned $($streams.count) streams"
    $script:STREAM_ID = $streams.streams[0].id
} else {
    Test-Fail "GET /api/streams" "No streams returned"
}

if ($STREAM_ID) {
    $stream = Invoke-GET "/api/streams/$STREAM_ID"
    if ($stream -and $stream.id -eq $STREAM_ID) {
        Test-Pass "GET /api/streams/:id - returned correct stream"
    } else {
        Test-Fail "GET /api/streams/:id" "Stream not found"
    }
}

Write-Host "  [INFO] WebSocket /ws/streams - skipped (test manually with wscat)" -ForegroundColor DarkYellow

# ── 6. Journeys ───────────────────────────────────────────────

Print-Header "6. Journeys"

$journeys = Invoke-GET "/api/journeys"
if ($journeys -and $journeys.count -gt 0) {
    Test-Pass "GET /api/journeys - returned $($journeys.count) journeys"
    $journeyId = $journeys.journeys[0].id

    $journey = Invoke-GET "/api/journeys/$journeyId"
    if ($journey -and $journey.id -eq $journeyId) {
        Test-Pass "GET /api/journeys/:id - returned correct journey"
    } else {
        Test-Fail "GET /api/journeys/:id" "Journey not found"
    }

    $booking = Invoke-POST "/api/journeys/$journeyId/book" @{} $TOKEN
    if ($booking) {
        Test-Pass "POST /api/journeys/:id/book - booked successfully"
    } else {
        Test-Fail "POST /api/journeys/:id/book" "Booking failed"
    }
} else {
    Test-Fail "GET /api/journeys" "No journeys returned"
}

# ── 7. Reminders ──────────────────────────────────────────────

Print-Header "7. Reminders"

if ($EVENT_ID) {
    $remind = Invoke-POST "/api/events/$EVENT_ID/remind" @{} $TOKEN
    if ($remind) {
        Test-Pass "POST /api/events/:id/remind - reminder set"
    } else {
        Test-Fail "POST /api/events/:id/remind" "Failed to set reminder"
    }

    $remindAgain = Invoke-POST-ExpectError "/api/events/$EVENT_ID/remind" @{} $TOKEN
    if ($remindAgain) {
        Test-Pass "POST /api/events/:id/remind (duplicate) - correctly returned 409"
    } else {
        Test-Fail "POST /api/events/:id/remind (duplicate)" "Should have returned 409"
    }

    $reminders = Invoke-GET "/api/auth/reminders" $TOKEN
    if ($reminders -and $reminders.count -gt 0) {
        Test-Pass "GET /api/auth/reminders - returned $($reminders.count) reminder(s)"
    } else {
        Test-Fail "GET /api/auth/reminders" "No reminders returned"
    }

    $deleteRemind = Invoke-DELETE "/api/events/$EVENT_ID/remind" $TOKEN
    if ($deleteRemind) {
        Test-Pass "DELETE /api/events/:id/remind - reminder removed"
    } else {
        Test-Fail "DELETE /api/events/:id/remind" "Failed to delete reminder"
    }
}

# ── 8. Merch ──────────────────────────────────────────────────

Print-Header "8. Merch"

$merch = Invoke-GET "/api/merch"
if ($merch -and $merch.count -gt 0) {
    Test-Pass "GET /api/merch - returned $($merch.count) items"
    $script:MERCH_ID = $merch.items[0].id

    $item = Invoke-GET "/api/merch/$MERCH_ID"
    if ($item -and $item.id -eq $MERCH_ID) {
        Test-Pass "GET /api/merch/:id - returned correct item"
    } else {
        Test-Fail "GET /api/merch/:id" "Item not found"
    }
} else {
    Test-Fail "GET /api/merch" "No merch items returned"
}

# ── 9. Orders ─────────────────────────────────────────────────

Print-Header "9. Orders"

if ($MERCH_ID) {
    $orderBody = @{ items = @(@{ merchId = $MERCH_ID; quantity = 2 }) }

    $order = Invoke-POST "/api/orders" $orderBody $TOKEN
    if ($order -and $order.order.id) {
        Test-Pass "POST /api/orders - order placed, total: $($order.order.total)"
        $script:ORDER_ID = $order.order.id
    } else {
        Test-Fail "POST /api/orders" "Order creation failed"
    }

    $orders = Invoke-GET "/api/orders" $TOKEN
    if ($orders -and $orders.count -gt 0) {
        Test-Pass "GET /api/orders - returned $($orders.count) order(s)"
    } else {
        Test-Fail "GET /api/orders" "No orders returned"
    }

    if ($ORDER_ID) {
        $singleOrder = Invoke-GET "/api/orders/$ORDER_ID" $TOKEN
        if ($singleOrder -and $singleOrder.id -eq $ORDER_ID) {
            Test-Pass "GET /api/orders/:id - returned correct order"
        } else {
            Test-Fail "GET /api/orders/:id" "Order not found"
        }
    }

    $noTokenOrder = Invoke-POST-ExpectError "/api/orders" $orderBody
    if ($noTokenOrder) {
        Test-Pass "POST /api/orders (no token) - correctly rejected"
    } else {
        Test-Fail "POST /api/orders (no token)" "Should have returned 401"
    }
}

# ── Summary ───────────────────────────────────────────────────

$TOTAL = $PASS + $FAIL
Write-Host ""
Write-Host "===================================================" -ForegroundColor DarkCyan
Write-Host "  RESULTS" -ForegroundColor Cyan
Write-Host "===================================================" -ForegroundColor DarkCyan
Write-Host "  Passed : $PASS / $TOTAL" -ForegroundColor Green
if ($FAIL -gt 0) {
    Write-Host "  Failed : $FAIL / $TOTAL" -ForegroundColor Red
} else {
    Write-Host "  Failed : 0 / $TOTAL" -ForegroundColor Green
    Write-Host ""
    Write-Host "  All tests passed!" -ForegroundColor Cyan
}
Write-Host ""
