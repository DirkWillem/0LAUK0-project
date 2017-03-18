import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { MedicationService, Medication } from "../core/services/medication.service";
import { PartialModel } from "../core/model";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";

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

  newMedication: PartialModel<Medication> = {
    title: "",
    description: ""
  };

  pendingMedication: Medication;

  errorMessage: string = "";

  constructor(private medicationService: MedicationService,
              private modalService: NgbModal,
              private route: ActivatedRoute,
              private router: Router) {

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

  /**
   * Event handler for opening the create medication modal
   * @param modal - The modal to open
   */
  openCreateMedicationModal(modal) {
    this.modalService.open(modal);
  }

  /**
   * Closes the create medication modal and clears the form
   * @param closeAction - The action to close the modal
   */
  closeCreateMedicationModal(closeAction) {
    closeAction();
    this.newMedication = { title: "", description: "" };
  }

  /**
   * Creates a new medication
   * @param closeAction - The action to close the creation modal
   * @returns {Promise<void>}
   */
  async createMedication(closeAction) {
    try {
      const medication = await this.medicationService.create(this.newMedication);
      this.router.navigate([medication.id], {relativeTo: this.route});
      closeAction();
    } catch(e) {
      this.errorMessage = e.message;
    }
  }

  /**
   * Starts the deletion of a medication
   * @param medication - The medication to delete
   * @param contents - The contents of the confirm delete modal
   */
  startDeleteMedication(medication: Medication, contents) {
    this.pendingMedication = medication;
    this.modalService.open(contents);
  }

  /**
   * Closes the confirm deletion modal
   * @param closeFn - The close function
   */
  closeConfirmDeleteModal(closeFn: () => void) {
    closeFn();
  }

  /**
   * Deletes the pending medication and closes the confirm medication modal
   * @param closeFn - The close function
   */
  async deletePendingMedication(closeFn: () => void) {
    await this.medicationService.delete(this.pendingMedication.id);
    this.medications = this.medications.filter(m => m.id != this.pendingMedication.id);
    this.closeConfirmDeleteModal(closeFn);
  }
}