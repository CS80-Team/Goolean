# Goolean

## Overview

<b>Goolean</b> is a boolean model search engine that allows users to perform boolean queries on a dataset of documents (corpus). The system consists of an engine that processes and indexes documents and a shell that interact with the engine through a set of commands.

## Features

- **Command-line Shell**: Provides an interface to interact with the engine.
- **Boolean Query Support**: Supports AND, OR, and NOT operations.
- **Document Indexing**: Efficiently indexes documents for fast retrieval.
- **Query Optimization**: Optimizes queries for better performance. (under construction)
- **File Management**: Allows adding new files and listing indexed documents.

## System Components

1. **Shell**: Accepts user commands and communicates with the engine.
3. **Engine**: Processes queries and retrieves relevant documents.
   - **Loader**: Loads documents into the system.
   - **Normalizer**: Prepares text by cleaning and standardizing it.
   - **Indexer**: Indexes documents for searching.
   - **Query Constructor**: Parses and builds the query.
   - **Query Optimizer**: Improves query execution efficiency. (Under construction)

## Commands

- `query`: Query the engine for a keyword or a boolean expression
  - Usage: `query <keyword> | <expression>`

- `list`: List all documents, displayable by <b>name</b> or/and <b>path</b> or/and <b>ID</b> or/and <b>extension</b>
Use -sortby to sort results by name, path, id or extension
Use -n to limit the number of results
Default fields order: `-id` `-name` `-path` `-ext`
Default sortby: `id`
Default limit: `all`
  - Usage: `list <-id | -name | -path | -ext> [-n <limit>] [-sortby <name | path | id | -ext>]`

- `load`: Load a new document into the engine
  - Usage: `load <document_path>`

- `find`: Find a document by name or id.
Display the document's id, name and path
Default search field: -name
  - Usage: `find <-id | -name> <value> || find <document_name>`

- `exit`: Exit the shell
  - Usage: `exit`

- `help`: List all available commands
  - Usage: `help`

- `clear`: Clear the screen
  - Usage: `clear`

- `open`: Open a document by ID in the default editor
  - Usage: `open <document_id>`

### Early Executed Commands (invoked after each command executed)

- `engine-stats`: Displays the total number of documents and keys in the engine.
  - Usage: `Reserved for the shell`

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/CS80-Team/Goolean/
   cd Goolean
   ```
2. Build the project:
   ```sh
   make build
   cd bin/
   ./goolean
   ```

## System Overview

![System Overview](https://github.com/CS80-Team/Goolean/blob/master/docs/BIRSystemOverview.png)

## Usage Example
[![Demo video](./docs/demo.gif)](./docs/demo.mp4)

## License

## Contributors

- [Omar Muhammad](https://github.com/OmarMGaber)
- [Ahmed Ashraf](https://github.com/ahmed-ashraff)
