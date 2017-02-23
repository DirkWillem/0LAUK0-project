import { Component, Input } from '@angular/core';

import { DoseSummarySummary, DoseStatus, DoseSummaryService } from "../core/services/dosesummary.service";

import * as moment from 'moment'

@Component({
  selector: "dose-summary",
  templateUrl: "./dosesummary.component.html",
  styleUrls: ["./dosesummary.component.scss"]
})
export class DoseSummaryComponent {
  @Input() summary: DoseSummarySummary;
  @Input() userId: number;
  opened: boolean = false;

  statuses: DoseStatus[] = null;

  constructor(private doseSummaryService: DoseSummaryService) {

  }

  /**
   * Formats a date string to a human readable string
   * @param date - The date to be formatted
   * @returns {string} - The formatted date
   */
  formatDate(date: string) {
    return moment(new Date(date)).format('MMMM Do YYYY');
  }

  /**
   * Event handler for toggling the item
   */
  async toggle() {
    if(!this.opened && this.statuses == null) {
      this.statuses = await this.doseSummaryService.listDoseStatuses(this.userId, this.summary.date);
    }

    this.opened = !this.opened;
  }
}