import { NgModule, Optional, SkipSelf } from '@angular/core';

import { DispatcherService } from "./services/dispatcher.service";
import { AuthService } from "./services/auth.service";
import { AuthHttp } from "./services/authhttp.service";
import { MedicationService } from './services/medication.service';
import { UserService } from "./services/user.service";
import { DoseService } from "./services/dose.service";
import { DoseSummaryService } from "./services/dosesummary.service";

import { MedicationsResolve, MedicationResolve } from './resolves/medication.resolve';
import { DoseSummariesResolve } from "./resolves/dosesummary.resolve";
import { PatientsResolve, UserResolve } from "./resolves/user.resolve";
import { DosesResolve } from "./resolves/dose.resolve";

import { AuthGuard} from './guards/auth.guard';
import { AutoLoginGuard } from "./guards/autologin.guard";
import { PRNMedicationService } from "./services/prnmedication.service";
import { PRNMedicationsResolve } from "./resolves/prnmedication.resolve";

/**
 * Module containing all services
 */
@NgModule({
  providers: [
    DispatcherService,
    AuthService, AuthGuard, AutoLoginGuard, AuthHttp,
    MedicationService, MedicationsResolve, MedicationResolve,
    UserService, PatientsResolve, UserResolve,
    DoseService, DosesResolve,
    DoseSummaryService, DoseSummariesResolve,
    PRNMedicationService, PRNMedicationsResolve
  ]
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parent: CoreModule) {
    if(parent) {
      throw new Error(
        'CoreModule is already loaded. Import it in the AppModule only');
    }
  }
}