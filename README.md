# Sack

*Sack* is a minimal command-line tool for creating web pages showcasing multiple 3D objects.

## Quick Start

Follow these steps to quickly install and run *Sack*:

1. Clone and download this repository.
2. In the top-level directory, run:

   ```bash
   go build -o sack ./cmd
   ```

3. Run the *Sack* program in your terminal. The default port is 7536.

   ```bash
   ./sack start
   ```

## Command Line Options

Use `./sack --help` for help. The main commands are:

- `start`: Starts the website at port 7536. Change the port with `--port`.
- `generate`: Generates a configuration list for your 3D object. When using `--batch` to generate multiple pages, ensure at least one-page configuration exists in `config.yaml` for the program to follow as an example. Specify the number of pages to generate after the `--batch` option.

## Project Structure

```plaintext
.
├── config.yaml
├── cmd/
│   ├── config.go
│   ├── config_test.go
│   ├── handlers.go
│   ├── helpers.go
│   ├── helpers_test.go
│   └── main.go
└── ui/
    ├── html/
    │   ├── 404.html
    │   ├── 500.html
    │   ├── index.html
    │   ├── pages/
    │   │   └── page1.gohtml
    │   └── templates/
    │       ├── base.gohtml
    │       ├── card.gohtml
    │       └── plain.gohtml
    └── static/
        ├── css/
        ├── img/
        ├── js/
        └── models/
            ├── example.glb
            ├── example.usdz
            └── example.webp
```

- The `cmd/` directory contains all the Go code and functions as a small server.
- The `ui/` directory is organized as follows:
  - `html/`: This directory holds HTML template files.
    - `index.html`: The main home page.
    - `pages/`: This subdirectory contains the generated HTML pages.
    - `templates/`: This subdirectory contains the template pages for generating HTML pages.
       - `base.gohtml`: The base template for all pages.
       - `card.gohtml`: Template for individual cards displaying 3D objects.
  - `static/`: This directory contains static files such as CSS, images, JavaScript, and 3D objects.
    - `css/`: Directory for CSS files.
    - `img/`: Directory for image files.
    - `js/`: Directory for JavaScript files.
    - `models/`: Directory for 3D model files. For optimal functionality, each 3D model should include:
      - `.usdz` for AR Quick Look on iOS devices.
      - `.glb` for efficient 3D rendering on the web.
      - `.webp` for high-quality, compressed poster images (`.png` and `.jpg` images also work).

The `config.yaml` file is the configuration file that generates multiple pages for our 3D objects. The required information format is:

```yaml
pages:
  page_name:
    ModelSrcPath: "/example/obj.glb"
    ModelIosSrcPath: "/example/obj.usdz"
    PosterPath: "/example/obj.webp"
    Description: "About_Me"
    ModelName: "Your_Model_Name"
    DesignerWebsite: "Your_Website"
    DesignerName: "Your_Name"
```

Add additional page configurations at the same level under `pages`.

## References

This work has largely benefited from the following projects:

1. [`<model-viewer>`](https://github.com/google/model-viewer)
2. [Bharat Icons](https://www.flaticon.com/authors/bharat-icons)

## Disclosure

The 3D object exhibited in this project is collected from [Dawanshiju](https://artsandculture.google.com/asset/aerial-view-of-dawanshiju/_QHjNn2iL_6JrQ?hl=en), Shenzhen, and is owned and shared by [Enza's Research Group](https://www.enzamigliore.com/).

## License

BSD-3. See [LICENSE](./LICENSE).
