/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./templ/*.templ",
               "./api.go" ],
    theme: {
      extend: {},
    },
    plugins: [],
  }

  // tailwind -i ./input.css -o ./templ/css/output.css --watch
  // tailwind -o ./templ/css/build.css --minify