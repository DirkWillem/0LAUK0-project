import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs'

import { DoseSummarySummary, DoseStatus, DoseSummaryService } from "../core/services/dosesummary.service";

import * as moment from 'moment'

@Component({
  selector: "dose-summary",
  templateUrl: "./dosesummary.component.html",
  styleUrls: ["./dosesummary.component.scss"]
})
export class DoseSummaryComponent implements OnInit, OnDestroy {
  @Input() summary: DoseSummarySummary;
  @Input() userId: number;
  opened: boolean = false;

  statuses: DoseStatus[] = null;
  statusesUpdatesSubscriptionId: number;
  statusesUpdatesSubscription: Subscription;

  constructor(private doseSummaryService: DoseSummaryService) {

  }

  /**
   * Initialization Angular lifecycle hook
   */
  async ngOnInit() {
    const statusesSubscription = await this.doseSummaryService.getDoseStatusesUpdates(this.userId, this.summary.date);
    this.statusesUpdatesSubscription = statusesSubscription.updates
      .subscribe(statuses => this.statuses = statuses);
    this.statusesUpdatesSubscriptionId = statusesSubscription.subscriptionId;
  }

  /**
   * Destruction Angular lifecycle hook
   */
  ngOnDestroy() {
    this.statusesUpdatesSubscription && this.statusesUpdatesSubscription.unsubscribe();
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
      console.log(this.summary);
      this.statuses = await this.doseSummaryService.listDoseStatuses(this.userId, moment(new Date(this.summary.date)).format("YYYY-MM-DD"));
    }

    this.opened = !this.opened;
  }
}