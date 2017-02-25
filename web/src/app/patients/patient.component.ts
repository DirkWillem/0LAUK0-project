import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { Subscription } from 'rxjs';

import { UserService, User } from "../core/services/user.service";
import { Dose, DoseService } from "../core/services/dose.service";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { DoseSummarySummary, DoseSummaryService } from "../core/services/dosesummary.service";
import { applyUpdateToCollection } from "../core/collectionupdates";

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
  doseSummaries: DoseSummarySummary[];
  routeDataSubscription: Subscription;

  pendingDose: Dose = null;

  constructor(private userService: UserService,
              private doseSummaryService: DoseSummaryService,
              private doseService: DoseService,
              private route: ActivatedRoute,
              private modalService: NgbModal) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  ngOnInit() {
    this.routeDataSubscription = this.route.data.subscribe(async (data: {patient: User, doses: Dose[], doseSummaries: DoseSummarySummary[]}) => {
      this.patient = data.patient;
      this.doses = data.doses;
      this.doseSummaries = data.doseSummaries;

      (await this.doseService.getCollectionUpdates(this.patient.id))
        .subscribe(mut => this.doses = applyUpdateToCollection(this.doses, mut));

      (await this.doseSummaryService.getDoseSummariesUpdates(this.patient.id))
        .subscribe(newSummaries => this.doseSummaries = newSummaries);
    });
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.routeDataSubscription && this.routeDataSubscription.unsubscribe();
  }

  /**
   * Event handler for opening the add dose modal
   * @param contents - The contents of the modal
   */
  openAddDoseModal(contents) {
    this.modalService.open(contents, {size: "lg"});
  }

  /**
   * Adds a created dose to the list of doses if it wasn't already added by the dispatcher
   * @param dose - The created dose
   */
  doseCreated(dose: Dose) {
    if(!this.doses.some(d => d.id == dose.id)) {
      this.doses.push(dose);
    }
  }

  /**
   * Event handler for opening the update dose modal
   * @param dose - The dose to update
   * @param contents - The contents of the update dose modal
   */
  openUpdateDoseModal(dose: Dose, contents) {
    this.pendingDose = dose;
    this.modalService.open(contents, {size: "lg"});
  }

  /**
   * Deletes a dose
   * @param dose - The dose to delete
   * @param event - The mouse event that triggered the delete
   * @returns {Promise<void>} - Promise that resolves once the delete is done
   */
  async deleteDose(dose: Dose, event: MouseEvent) {
    event.stopPropagation();
    event.preventDefault();

    await this.doseService.delete(this.patient.id, dose.id);
    this.doses = this.doses.filter(d => dose.id != d.id);
  }

  /**
   * Updates the doses list after a dose has been updated
   * @param updatedDose - The updated dose
   */
  doseUpdated(updatedDose: Dose) {
    this.doses = this.doses.map(dose => dose.id == updatedDose.id ? updatedDose : dose);
  }
}