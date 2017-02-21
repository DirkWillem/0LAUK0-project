import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';

import { MedicationService, Medication } from "../core/services/medication.service";

/**
 * Medication list component
 */
@Component({
  selector: "medications",
  templateUrl: "./medications.component.html"
})
export class MedicationsComponent implements OnInit, OnDestroy {
  medications: Medication[];
  routeDataSubscription: Subscription;

  constructor(private medicationService: MedicationService, private route: ActivatedRoute) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  ngOnInit() {
    this.routeDataSubscription = this.route.data.subscribe((data: {medications: Medication[]}) => {
      this.medications = data.medications;
    });
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.routeDataSubscription && this.routeDataSubscription.unsubscribe();
  }
}