import { Component, OnInit } from '@angular/core';
import { addDays, startOfWeek } from 'date-fns';
import { Session } from 'src/app/models/session.model';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss']
})
export class CalendarComponent implements OnInit {

  days: Date[] = [];
  sessions: Session[] = [{
    day: new Date(),
    name: 'Schwimmen'
  }];

  constructor() { }

  ngOnInit() {
    const today = new Date();
    const monday = startOfWeek(today, { weekStartsOn: 1 });


    for (let i = 0; i < 7; i++) {
      this.days.push(addDays(monday, i));
    }
  }

}
