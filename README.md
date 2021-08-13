# **Valinar Kinetic**
Repository for Valinar kinetic materials.

## **Scripts**
Contains Icarus scripts for drone operation.

## **TEMPLATEroutes.xlsx**
Table for mass drone movement.

Columns:
* Launch - Time in seconds to wait before sending drone
* Lat/Lon - Coordinates to go to
* Alt/Vel - Target altitude and speed
* Action - What action to perform for the waypoint
  * GOTO (_blank_) - Go to waypoint
  * LINGER - Go to waypoint and wait
    * Must input number of seconds to wait as parameter
  * FIRE - Go to waypoint and wait--shoot on sight
    * Must input number of seconds to wait as parameter
* Param - Parameter for the action

### **UI**
Contains UI files and compatible Icarus scripts.

To build all go files into executables
* `./build.sh`

#### **Dependencies**
To use the UI, gtk2 is needed  
 * `sudo pacman -S gtk2`

To build the UI, glade is needed
 * `sudo pacman -S glade`

### **Reference**
Contains outdated kinetic materials and example scripts to be used as a reference.

## **Delphi2**
Contains files for Delphi2.


# ACE-2021-Capstone-Icarus
# ACE-2021-Capstone-Icarus
