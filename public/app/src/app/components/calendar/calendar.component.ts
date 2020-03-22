import { Component, OnInit, ViewChild } from '@angular/core';
import { addDays, startOfWeek } from 'date-fns';
import { Session } from 'src/app/models/session.model';
import { ModalComponent } from '../modal/modal.component';
import { SessionService } from 'src/app/services/session.service';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss']
})
export class CalendarComponent implements OnInit {
  @ViewChild('modal', { static: false }) modal: ModalComponent;

  days: Date[] = [];
  selectedSession: Session;
  sessions: Session[] = [];

  constructor(private sessionService: SessionService) { }

  ngOnInit() {
    const today = new Date();
    const monday = startOfWeek(today, { weekStartsOn: 1 });

    // build the current week
    for (let i = 0; i < 7; i++) {
      this.days.push(addDays(monday, i));
    }

    // gets the sessions
    this.sessionService.data$.subscribe(sessions => {
      this.sessions = [...sessions];
    });
  }

  /**
   * Calculates the top position, based on the time of the event
   */
  getTopPosition(session: Session): number {
    const start = Number(session.startTime.replace(':', ''));

    return (start - 1200) * 0.5;
  }

  /**
   * Calculates the height, based on the duration of the event
   */
  getHeight(session: Session): number {

    const start = Number(session.startTime.replace(':', ''));
    const end = Number(session.endTime.replace(':', ''));

    const duration = end - start;
    return duration * 0.8;


  }

  /**
   * Open session details modal
   */
  openSessionDetails(session: Session): void {
    this.modal.show();
    this.selectedSession = session;
  }

}
