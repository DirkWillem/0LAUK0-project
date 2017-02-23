import { Injectable } from '@angular/core';

import { AuthHttp } from './authhttp.service';
import { Model, Field, ModelJson } from "../model";

/**
 * Contains summary information on a dose summary
 */
export class DoseSummarySummary extends Model {
  @Field() date: string;
  @Field() dispensedCount: number;
  @Field() pendingCount: number;
  @Field() totalCount: number;
}

/**
 * Contains the status of a dose on a given date
 */
export class DoseStatus extends Model {
  @Field() dispensedTime: string;
  @Field() dispensed: boolean;
  @Field() pending: boolean;
  @Field() beingDispensed: boolean;
  @Field() dose: {id: number, title: string};
}

/**
 * Service for interfacing with the dose summary API
 */
@Injectable()
export class DoseSummaryService {
  constructor(private http: AuthHttp) {

  }

  /**
   * Returns all dose summaries for a given user ID
   * @param userId - The ID of the user to list the dose summaries for
   * @returns {Promise<DoseSummarySummary[]>} - Promise resolving to the dose summaries
   */
  async listDoseSummaries(userId: number): Promise<DoseSummarySummary[]> {
    const json = await this.http.getJSON<ModelJson<DoseSummarySummary>[]>(`/api/users/${userId}/dosesummaries`);
    return json.map(item => new DoseSummarySummary(item));
  }

  /**
   * Returns all dose statuses for a given user ID on a given date
   * @param userId - The ID of the user the dose statuses belong to
   * @param date - The date to find the dose statuses for
   * @returns {Promise<DoseStatus[]>} Promise resolving to the dose statuses
   */
  async listDoseStatuses(userId: number, date: string): Promise<DoseStatus[]> {
    const json = await this.http.getJSON<ModelJson<DoseStatus>[]>(`/api/users/${userId}/dosesummaries/${date}`);
    return json.map(item => new DoseStatus(item));
  }
}