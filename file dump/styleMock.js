module.exports = new Proxy(
  {},
  {
    get: function (_, prop) {
      return prop
    },
  }
)
