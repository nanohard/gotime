# gotime
Time tracker for projects and tasks written in golang.

[![GoDoc](https://godoc.org/github.com/nanohard/gotime?status.svg)](https://godoc.org/github.com/nanohard/gotime)

**Only supported on Linux**

![gotime](https://user-images.githubusercontent.com/10169206/27987771-f1a5bd9e-63d8-11e7-8d3d-a8abc33bd0e9.gif)

## Overview
I created this in response to my need to record entries for tasks and tasks for projects, and my want to not have to type a whole command in the console for every action I wanted to perform. Thus gotime was born; an ncurses-style console user interface program with an SQlite3 database.

The time is stamped in seconds upon starting and stopping an entry; no timer is actually used. Time is rounded to the nearest minute and is viewable for each entry, task, and project alongside each item's description.

## Getting Started
Upon starting the program you will be required to type in the name of a project. Navigate to the right and press Ctrl-A to add a task, and do the same to create and start an entry. Press Ctrl-S to save the entry.

## Controls
* **Ctrl-A**: Add an item to the current view. If in the entries view it will "start" an Entry. This will create a timestamp of the current time and create an entry in the database. At this point you may type notes into the Entry's details. Press Ctrl-S to stop and save the Entry.
* **Ctrl-S**: Save the text you have written in the output box, whether for an Entry's details or a description for a Project or Task. For an Entry this will save the details you have written and stop the timer.
* **Ctrl-D**: Add a description to a Project or Task.
* **Ctrl-R**: Remove an Entry, Task, or Project. Removing a Project or Task will also remove all of its children.
* **Arrow Keys**: Left and right will navigate between Projects, Tasks, and Entries. Up and down will navigate within projects, tasks, and entries.
* **Ctrl-C**: Quit the program.

## ToDo
- [ ] Windows compatible
- [ ] Mac compatible
- [ ] Integrate Harvest API
- [ ] Allow for archiving


## Credits
Created by [**Nanohard**](https://github.com/nanohard)

### Libraries
* [gocui](https://github.com/jroimartin/gocui)
* [gorm](https://github.com/jinzhu/gorm)

### License
Released under the [BSD-3 License](https://github.com/nanohard/gotime/blob/master/LICENSE)
