# Cat API Viewer

This project is a web application that integrates with The Cat API to display cat images and information. It uses Beego for the backend and serves a pre-built React frontend through Beego's template rendering system.

## Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Generating an API Key](#generating-an-api-key)
5. [Configuration](#configuration)
6. [Running the Application](#running-the-application)
7. [Project Structure](#project-structure)
8. [Development](#development)
9. [Frontend Development (Optional)](#frontend-development-optional)
10. [API Calls](#api-calls)
11. [Contributing](#contributing)
12. [License](#license)

## Features

- Browse cat images from The Cat API
- Filter cats by breed
- View detailed information about each cat
- Add cats to favorites
- View all favorite cats
- Vote and unvote for cats
- Responsive design for various screen sizes

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (version 1.16 or later)
- Node.js (version 14 or later)
- npm (usually comes with Node.js)
- Git

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/noman1811048/catApp_backend.git
   cd catApp_backend
   ```

2. **Install Dependencies:**
   Ensure you have Go installed, along with the Beego CLI tool (`bee`). If you don't have `bee` installed, you can do so with:

   ```bash
   go install github.com/beego/bee/v2@latest
   ```

   Then, install project dependencies:

   ```bash
   go mod tidy
   ```

3. **Set up your configuration:**
   Create a `conf/app.conf` file in your project root with the following content:

   ```
   appname = catApi
   httpport = PORT ADDRESS
   runmode = dev
   StaticDir = static:static
   cat_api_key = YOUR_CAT_API_KEY_HERE
   ```

   Replace `YOUR_CAT_API_KEY_HERE` with the API key you received from The Cat API.

Note: The React project has already been built and integrated into the Beego project structure. The built files like JS and CSS are located in the `static` folder, and the `index.html` has been renamed to `index.tpl` in the `views` folder.

## Generating an API Key

To use The Cat API, you need to generate an API key:

1. Visit <https://thecatapi.com/>
2. Click on the "GET YOUR API KEY" button
3. Fill out the registration form with your email address
4. Check your email for a message from The Cat API containing your API key
5. Copy this API key and add it to your `conf/app.conf` file

## Configuration

Ensure that your `conf/app.conf` file contains the necessary configuration parameters, including your Cat API key, as shown in the [Installation](#installation) section.

## Running the Application

1. Start the Beego server:

   ```bash
   bee run
   ```

2. Open your browser and navigate to `http://localhost:8080/`

Note: There's no need to run the frontend separately as it's already built and integrated with the Beego project for template rendering.

## Project Structure

```
catApp_backend/
├── conf/
│   ├── app.conf
│   └── app.conf.sample
├── controllers/
│   ├── cat_controller.go
│   └── default.go
├── routers/
│   └── router.go
├── static/
│   ├── css/
│   └── js/
├── tests/
│   └── default_test.go
├── views/
│   └── index.tpl
├── .gitignore
├── catApp_backend
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Development

For backend development:

1. Make changes to Go files in the appropriate directories
2. Run `bee run` to start the Beego server with hot reloading

## Frontend Development (Optional)

If you need to make changes to the frontend:

1. Clone the original React project:

   ```bash
   git clone https://github.com/noman1811048/catApp.git
   cd catApp
   ```

2. Install dependencies: `npm install`
3. Start the development server: `npm run dev`
4. Make your changes in the React app
5. Build the app when ready: `npm run build`
6. Copy the built files to the appropriate location in the Beego project structure:
   - Copy all files from the `dist` folder to the `static` folder in the Beego project
   - Rename `index.html` to `index.tpl` and move it to the `views` folder in the Beego project

Note: This step is only necessary if you need to modify the frontend. The current setup already includes a built version of the frontend.

## API Calls

The application uses Go channels for making API calls to The Cat API. This allows for efficient, non-blocking operations when fetching cat data.

## Contributing

Contributions to this project are welcome. Please fork the repository and submit a pull request with your changes.

## License

MIT License
Copyright (c) 2024 Asadullah Al Noman
