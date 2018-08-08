package database

import (

)

type ScheduledEvent struct {
	ScheduleName 	string	`db:"schedulename"`
	EName			string  `db:"ename"` 
	//details
	EReceiveBy		string  `db:"ereceiveby"` 


}