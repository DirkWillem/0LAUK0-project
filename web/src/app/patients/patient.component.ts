import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { Subscription } from 'rxjs';

import { UserService, User } from "../core/services/user.service";
import { Dose } from "../core/services/dose.service";

/**
 * Component for a single patient
 */
@Component({
  selector: "patient",
  templateUrl: "./patient.component.html"
})
export class PatientComponent implements OnInit, OnDestroy {
  patient: User;
  doses: Dose[];
  routeDataSubscription: Subscription;

  constructor(private userService: UserService, private route: ActivatedRoute) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  ngOnInit() {
    this.routeDataSubscription = this.route.data.subscribe((data: {patient: User, doses: Dose[]}) => {
      this.patient = data.patient;
      this.doses = data.doses;
    });
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.routeDataSubscription && this.routeDataSubscription.unsubscribe();
  }
}