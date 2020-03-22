import { Component, OnInit, ViewChild } from '@angular/core';
import { AuthService } from 'src/app/auth/auth.service';
import { SessionService } from 'src/app/services/session.service';
import { ModalComponent } from 'src/app/components/modal/modal.component';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Session } from 'src/app/models/session.model';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  @ViewChild('modal', { static: false }) modal: ModalComponent;

  sessionForm: FormGroup;

  constructor(public auth: AuthService,
              private formBuilder: FormBuilder,
              private sessionService: SessionService) { }

  ngOnInit() {
    this.sessionForm = this.formBuilder.group({
      name: [''],
      description: [''],
      date: ['20.03.2020'],
      startTime: ['16:00'],
      endTime: ['17:00']
    });
  }

  /**
   * Scroll button event handler
   */
  scrollToCalendar(): void {
    window.scrollBy({
      top: 600,
      left: 0,
      behavior: 'smooth'
    });
  }


  /**
   * Submit form to add new session
   */
  addNewSession(): void {
    const newSession = new Session({
      name: this.sessionForm.get('name').value,
      description: this.sessionForm.get('description').value,
      date: new Date(this.sessionForm.get('date').value),
      startTime: this.sessionForm.get('startTime').value,
      endTime: this.sessionForm.get('endTime').value
    });
    this.sessionService.addSession(newSession);

    this.modal.hide();
    this.sessionForm.reset();
  }

}
