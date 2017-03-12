import { NgModule } from '@angular/core';
import { FormatDatePipe } from "./pipes/dateformat.pipe";

@NgModule({
  declarations: [FormatDatePipe],
  exports: [FormatDatePipe]
})
export class SharedModule {

}