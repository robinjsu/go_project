### Building a Concurrent GUI in GO - Special Project (IN PROGRESS AS OF APRIL 2022)

##### Author: Robin Su
##### Advisor: Andrew Tolmach
##### Portland State University
##### Spring 2022 

Description: This project centers around building a web-based application to explore the Go programming language, with a focus on concurrency management and multi-threading. The foundation of this project is based on the design by [Michal Štrba](https://github.com/faiface/gui), which in turn was inspired by Rob Pike's ideas on concurrency and message passing in [Newsqueak](https://www.youtube.com/watch?v=hB05UFqOtFA&t=2408s).

This project includes a comparative study of Go and Python, in order to understand how the two languages differ in achieving similar functionality.

The main architecture of this application separates the GUI into components, each of which runs as a separate goroutine. Communication between components in achieved through channels, a concurrency primitive built into the Go programming language. Rather than using shared memory to signal changes in state, message passing is the main method to do so.

-------

#### Planned Application Design:

![gotextaide](https://user-images.githubusercontent.com/53282793/164770389-36072824-009b-45ef-b026-8fad29d952c0.jpg)

-------

### Additional Resources:

 - K. Cox-Buday *Concurrency in Go: Tools & Techniques for Developers* Sebastopol: O'Reilly Media, Inc. 2017.

 - N. Tao "The Go image/draw package - The Go Programming Language". Blog: https://go.dev/blog/image-draw. [Accessed: March/April 2022]
