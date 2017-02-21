import { NgModule } from '@angular/core';
import { BrowserModule } from "@angular/platform-browser";
import { FormsModule } from "@angular/forms";
import { HttpModule } from "@angular/http";

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { appComponents } from "./app.declarations";
import { AppComponent } from "./app.component";
import { CoreModule } from "./core/core.module";
import { appRoutingModule } from "./app.routing";

/**
 * Main application module
 */
@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    CoreModule,
    NgbModule.forRoot(),
    appRoutingModule
  ],

  providers: [],
  declarations: [...appComponents],
  bootstrap: [AppComponent]
})
export class AppModule {

}