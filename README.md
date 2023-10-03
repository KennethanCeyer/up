<h1 align="center">up</h1>
<p align="center">Experimental: Thread-friendly interpreter programming language</p>
<p align="center"><img src="./logo.svg" width="140" alt="up" /></p>

## Features

- Optimized for lightweight threading.
- Strongly influenced by both the Go (Golang) and Python programming languages.
- While garbage collection (GC) exists, it minimizes lock overhead by leveraging the CAS Atomic Operator.
- It's a dynamically typed language; however, users can specify strong types for instances where high performance is a necessity.
- Features a straightforward and easily comprehensible grammatical structure.
- Developed as an experimental language, intended for long-term project exploration and refinement.

## Overview

The `up` language aspires to marry the ease of use seen in languages like Python with the concurrency advantages inherent to languages like Go. Its design philosophy centers on offering a platform where developers can utilize threads without confronting the usual associated intricacies. As a unique proposition in the programming language spectrum, `up` merits exploration by those keen on advancing the paradigms of contemporary programming.

## Getting Started

```bash
go run . {.up file path}
```

Example:

```bash
go run . examples/for_loop.up
```

This project is working in progress. There may be an error in the code's behavior.

## Contributing

Contributions to this project are appreciated. For additional details, please consult the CONTRIBUTING.md document.

## License

This project is licensed under the Apache License 2.0. Detailed licensing information can be accessed in the [LICENSE](LICENSE) file.
