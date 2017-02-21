import { NgModule, Optional, SkipSelf } from '@angular/core';

import { AuthService } from "./services/auth.service";
import { AuthHttp } from "./services/authhttp.service";
import { APIInterface } from './services/apiinterface.service';

import { MedicationService } from './services/medication.service';

import { MedicationsResolve, MedicationResolve } from './resolves/medication.resolve';

import { AuthGuard} from './guards/auth.guard';
import { AutoLoginGuard } from "./guards/autologin.guard";
import { UserService } from "./services/user.service";

/**
 * Module containing all services
 */
@NgModule({
  providers: [
    AuthService, AuthGuard, AutoLoginGuard, AuthHttp,
    APIInterface,
    MedicationService, MedicationsResolve, MedicationResolve,
    UserService
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