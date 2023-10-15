# Project LEM-IN

Authors Toomas Vooglaid (tvooglai) and Kurban Ramazanov

Task [description](https://github.com/01-edu/public/tree/master/subjects/lem-in) and audit [requirements](https://github.com/01-edu/public/tree/master/subjects/lem-in/audit)

## Usage

1. Clone repository: `git clone https://01.kood.tech/git/tvooglai/lem-in.git`
2. Move into folder: `cd lem-in`
3. Run visualisation by executing `./run.sh [example01.txt]` or follow steps manually:
    * To get ants' steps textually, execute `go run . example01.txt`
    * Change number of ants by providing option `-ants <int>`, e.g. `go run . -ants 5 example01.txt`
    * Run visualisation by executing `go run . example01.txt | go run .`
        * Adjust terminal size to accommodate all rooms and paths
        * Follow instructions in bottom of page
        * Change size of graph by providing zooming options to second program, e.g. `go run . example01.txt | go run . -x 5 -y 3`
    * Prepare binaries by executing `./prepare.sh` or prepare binaries yourself by following example in this file
    * Run program `lem-in` alone to see textual information, or pipe the output to `visualizer` to see visualisation, e.g. `./lem-in -ants 20 example05.txt | ./visualizer -x 5 -y 2`

### Run browser visualisation:

You can run browser visualisation by providing `-web` option, e.g. `go run . -web example05.txt`.

Open in your browser: `http://localhost:8080/`.

On the page, you will find an intuitive interface where you can find all possible paths in the graph and view the visualization of ant movement. You can also choose an ant for path traversal.


## Comment about direction of ants' movements

As ants stubbornly face left side (see e.g. visualisation of example00), we had to change start and end points to match their preferences :).

### Disclaimer

Initial inspiration for this project was sparked from Egon Saks [work](https://github.com/EgonSaks/lem-in).