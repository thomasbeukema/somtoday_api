#!/usr/bin/env python3

"""
	Author: Thomas Beukema
	Email: thomasbeukema@outlook.com

	This file is an example on how to use the Python version of the SOMToday API
"""

from somtoday import * # import all function from our API

"""
	Create tuple with your credentials stored in it.
		0 => This is your username you use to log in to SOMToday
		1 => This is your password you use to log in to SOMToday
		2 => This is the abbreviation of your school*.
		3 => This is the BRIN code of your school*.

		* You can find all this on servers.somtoday.nl
"""
credentials = ( "YOUR USERNAME", "YOUR PASSWORD", "YOUR SCHOOL", "BRIN CODE OF YOUR SCHOOL" )

som = Somtoday(credentials) # Instantiate the API class

timetable = som.getTimetable() # Get your timetable
homework  = som.getHW()		   # Get your homework
grades    = som.getGrades()	   # Get your grades

# Change the status of your homework item.
# This function returns True when it has success, otherwise it will return False.
if som.changeHwDone(1073870857, 1089708218, False):
	print("Changed the status of the homework item successfully")
else:
	print("Failed to change the status of the homework item")