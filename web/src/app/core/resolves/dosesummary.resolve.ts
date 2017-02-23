import { Injectable } from "@angular/core";
import { Resolve, ActivatedRouteSnapshot } from '@angular/router';

import { DoseSummaryService, DoseSummarySummary } from '../services/dosesummary.service';

/**
 * Resolve for a list of dose summaries
 */
@Injectable()
export class DoseSummariesResolve implements Resolve<DoseSummarySummary[]> {
  constructor(private doseSummaryService: DoseSummaryService) {

  }

  async resolve(route: ActivatedRouteSnapshot): Promise<DoseSummarySummary[]> {
    let r = route;
    while(!r.params["id"]) {
      r = r.parent;
    }

    const id = r.params["id"];

    return await this.doseSummaryService.listDoseSummaries(parseInt(id));
  }
}