import { Component, OnInit, ViewChild } from '@angular/core';
import { addDays, startOfWeek } from 'date-fns';
import { Session } from 'src/app/models/session.model';
import { ModalComponent } from '../modal/modal.component';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss']
})
export class CalendarComponent implements OnInit {
  @ViewChild('modal', { static: false }) modal: ModalComponent;

  days: Date[] = [];
  selectedSession: Session;
  sessions: Session[] = [
    {
      day: new Date(),
      name: 'Schwimmen'
    },
    {
      day: addDays(new Date(), -2),
      name: 'Lachen'
    },
    {
      day: addDays(new Date(), -4),
      name: 'Tanzen'
    }
  ];

  constructor() { }

  ngOnInit() {
    const today = new Date();
    const monday = startOfWeek(today, { weekStartsOn: 1 });

    // build the current week
    for (let i = 0; i < 7; i++) {
      this.days.push(addDays(monday, i));
    }
  }

  /**
   * Calculates the top position, based on the time of the event
   */
  getTopPosition(session: Session): number {
    // TODO: fake data
    return Math.random() * 560;
  }

  /**
   * Calculates the height, based on the duration of the event
   */
  getHeight(session: Session): number {
    // TODO: fake data
    return Math.random() * 160 + 40;
  }

  /**
   * Open session details modal
   */
  openSessionDetails(session: Session): void {
    this.modal.show();
    this.selectedSession = session;
  }

}
