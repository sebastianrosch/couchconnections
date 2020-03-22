import { Pipe, PipeTransform } from '@angular/core';
import { isSameDay } from 'date-fns';
import { Session } from '../models/session.model';

@Pipe({
  name: 'sessionsOnDay'
})
export class SessionsOnDayPipe implements PipeTransform {

  transform(array: Session[], day: Date): Session[] {

    return array.filter(session => isSameDay(session.date, day));
  }

}
