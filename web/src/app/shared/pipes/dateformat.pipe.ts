import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';

@Pipe({
  name: "formatDate"
})
export class FormatDatePipe implements PipeTransform {
  transform(date: Date): string {
    return moment(date).format("MMMM Do YYYY");
  }
}