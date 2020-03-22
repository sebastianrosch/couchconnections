import { Injectable } from '@angular/core';
import { Session } from '../models/session.model';
import { sessions } from './sessions.data';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SessionService {

  private data = sessions;
  data$ = new BehaviorSubject<Session[]>(this.data);

  constructor() {
  }

  /**
   * Adds a new session
   */
  addSession(newSession: Session): void {
    this.data.push(newSession);
    this.data$.next(this.data);
    console.log('added');
  }

}
