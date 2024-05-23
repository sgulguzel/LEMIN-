OVERVIEW 
lem-in is a Go-based program that simulates an ant farm. The goal is to move a specified number of ants from a start room (##start) to an end room (##end) in the shortest time possible. The program reads the colony configuration from a file provided as a command-line argument, validates the input, and then simulates the movement of ants across the colony. The output includes the original input data and the sequence of moves made by the ants.
 

FEATURE
Reads colony configuration from a file.
Validates the format and content of the input.
Finds all possible paths from the start to the end room.
Filters paths to avoid overlap and maximize efficiency.
Simulates the movement of ants and prints each move.


HOW TO RUN IT
The program initially reads a colony configuration file, which contains the specifications of ants and the colony. Then, it employs a path-finding algorithm to determine all possible paths from the start room to the end room. It selects the most suitable paths from these options and filters them to prevent collisions. Subsequently, it simulates the movement of ants along these filtered paths, ensuring proper relocation of ants at each turn and tracking their efforts to reach the destination rooms. Upon completion of these processes, the program generates a detailed output depicting the colony's configuration and the movement of ants.

USAGE
go run main.go <file_name>

AUTHOR
SÜMEYYE GÜL GÜZEL