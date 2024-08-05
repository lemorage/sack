# Sack

*Sack* is a minimal command-line tool for creating web pages that showcase multiple 3D objects.

## Quick Start

Follow these steps to quickly install and run *Sack*:

1. **Clone and download this repository.**

2. **Run the setup script:**
   ```sh
   chmod +x setup.sh
   ./setup.sh
   ```

3. **Run the application:**
   ```sh
   ./sack start
   ```

The default port is 7536. You can change the port using command-line options.

## Command Line Options

Use `./sack --help` for help. The main commands are:

- `start`: Starts the website on port 7536. Change the port with `--port`.
- `generate`: Generates a configuration list for your 3D objects. Use `--batch` to generate multiple pages, ensuring at least one-page configuration exists in `config.yaml` as an example. Specify the number of pages to generate after the `--batch` option.

## Project Structure

```plaintext
.
├── configs/
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
    │   ├── graph.html
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

- `configs/`: Configuration files for the application.
- `cmd/`: Contains the Go code and functions as a small server.
- `ui/`: Contains HTML templates and static files.
  - `html/`: Holds HTML template files.
    - `index.html`: The main home page.
    - `graph/html`: The story graph page.
    - `pages/`: Contains generated HTML pages.
    - `templates/`: Templates for generating HTML pages.
      - `base.gohtml`: The base template for all pages.
      - `card.gohtml`: Template for individual cards displaying 3D objects.
  - `static/`: Contains static files such as CSS, images, JavaScript, and 3D objects.
    - `css/`: CSS files.
    - `img/`: Image files.
    - `js/`: JavaScript files.
    - `models/`: 3D model files, including:
      - `.usdz` for AR Quick Look on iOS devices.
      - `.glb` for efficient 3D rendering on the web.
      - `.webp` for high-quality, compressed poster images (also supports `.png` and `.jpg`).

## Configuration File

### 1. `config.yaml`

The `config.yaml` file generates multiple pages for your 3D objects. Use the following format:

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

### 2. `graph.json`

The `graph.json` file generates the story graph based on nodes and links. Use the following format:

```json
{
  "nodes": [
    {
      "id": "node1",
      "keyword": "citrus",
      "story": "I am sweet"
    },
    {
      "id": "node2",
      "keyword": "pomelo",
      "story": "I am bitter"
    }
  ],
  "links": [
    {
      "source": "node1",
      "target": "node2"
    }
  ]
}
```

Add additional nodes and links as needed to build your story graph.

## References

This project has benefited from the following:

1. [`<model-viewer>`](https://github.com/google/model-viewer)
2. [`three.js`](https://threejs.org)
3. [`D3.js`](https://d3js.org)
4. [Bharat Icons](https://www.flaticon.com/authors/bharat-icons)

## Disclosure

The 3D stone exhibited in this project is collected from [Dawanshiju](https://artsandculture.google.com/asset/aerial-view-of-dawanshiju/_QHjNn2iL_6JrQ?hl=en), Shenzhen, and is kindly shared by [Enza's Research Group](https://www.enzamigliore.com/).

## License

BSD-3. See [LICENSE](./LICENSE).
