/**
 * Enum containing all possible field types
 */
enum FieldType {
  PRIMITIVE = 0,
  DATE = 2,
  MODEL = 3,
  MODEL_LIST = 4,
  MODEL_DICT = 5
}

/**
 * Interface for model field metadata
 */
interface ModelField {
  name: string;
  type: FieldType;
  model?: ModelClass;
  detail: boolean;
  optional: boolean;
}

/**
 * Type for model classes
 */
export type ModelClass = {new(json, stub?: boolean): Model};
export type SubModelClass<T extends Model> = {new(json, stub?: boolean): T}

/**
 * Error class for model errors
 */
class ModelError extends Error {

}

/**
 * Decorator for a plain-value field
 * @param detail - Whether the field is only present in the detail version of the model
 * @param optional - Whether the field is optional
 */
export function Field({detail, optional} :
  {detail?: boolean, optional?: boolean} = {detail: true, optional: false}) {
  return (target, fieldName) => {
    if(target.fields instanceof Array) {
      target.fields.push({name: fieldName, type: FieldType.PRIMITIVE, detail, optional});
    } else {
      target.fields = [{name: fieldName, type: FieldType.PRIMITIVE, detail, optional}];
    }
  }
}

/**
 * Decorator for a date field
 * @param detail - Whether the field is only present in the detail version of the model
 * @param optional - Whether the field is optional
 */
export function DateField({detail, optional} :
  {detail?: boolean, optional?: boolean} = {detail: true, optional: false}) {
  return (target, fieldName) => {
    if(target.fields instanceof Array) {
      target.fields.push({name: fieldName, type: FieldType.DATE, detail, optional});
    } else {
      target.fields = [{name: fieldName, type: FieldType.DATE, detail, optional}];
    }
  }
}

/**
 * Decorator for a field containing a sub-model
 * @param detail - Whether the field is only present in the detail version of the model
 * @param optional - Whether the field is optional
 * @param model - The model class of the sub-model
 */
export function ModelField({detail, optional, model} :
  {detail?: boolean, optional?: boolean, model: {new(json): Model}} = {detail: true, optional: false, model: Model}) {
  return (target, fieldName) => {
    if(target.fields instanceof Array) {
      target.fields.push({name: fieldName, type: FieldType.MODEL, model, detail, optional});
    } else {
      target.fields = [{name: fieldName, type: FieldType.MODEL, model, detail, optional}];
    }
  }
}

/**
 * Decorator for a field containing a list of sub-models
 * @param detail - Whether the field is only present in the detail version of the model
 * @param optional - Whether the field is optional
 * @param model - The model class of the sub-model in the list
 */
export function ModelListField({detail, optional, model} :
  {detail?: boolean, optional?: boolean, model: {new(json): Model}} = {detail: true, optional: false, model: Model}) {
  return (target, fieldName) => {
    if(target.fields instanceof Array) {
      target.fields.push({name: fieldName, type: FieldType.MODEL_LIST, model, detail, optional});
    } else {
      target.fields = [{name: fieldName, type: FieldType.MODEL_LIST, model, detail, optional}];
    }
  }
}

/**
 * Decorator for a field containing a dictionary of sub-models
 * @param detail - Whether the field is only present in the detail version of the model
 * @param optional - Whether the field is optional
 * @param model - The model class of the sub-model in the dictionary
 */
export function ModelDictField({detail, optional, model} :
  {detail?: boolean, optional?: boolean, model: {new(json): Model}} = {detail: true, optional: false, model: Model}) {
  return (target, fieldName) => {
    if(target.fields instanceof Array) {
      target.fields.push({name: fieldName, type: FieldType.MODEL_DICT, model, detail, optional});
    } else {
      target.fields = [{name: fieldName, type: FieldType.MODEL_DICT, model, detail, optional}];
    }
  }
}

/**
 * Represents partial model data
 */
export type PartialModel<M extends Model> = {
  [P in keyof M]?: M[P];
}

export type ModelJson<M extends Model> = {
  readonly [P in keyof M]: M[P];
}

/**
 * Base class for all models
 */
export class Model {
  private static fields: ModelField[];

  private get modelFields() {
    return this.constructor.prototype.fields;
  }

  constructor(json, isStub?: boolean) {
    if(!isStub) {
      for(let field of this.modelFields) {
        const fieldName = field.name;
        if(json.hasOwnProperty(fieldName)) {
          switch(field.type) {
            case FieldType.PRIMITIVE:
              this[fieldName] = json[fieldName];
              break;
            case FieldType.DATE:
              this[fieldName] = new Date(json[fieldName]);
              break;
            case FieldType.MODEL:
              this[fieldName] = new field.model(json[fieldName]);
              break;
            case FieldType.MODEL_LIST:
              this[fieldName] = json[fieldName].map(item => new field.model(item));
              break;
            case FieldType.MODEL_DICT: {
              let obj = {};
              for(let key in json[fieldName]) {
                obj[key] = new field.model(json[fieldName][key]);
              }
              this[fieldName] = obj;
              break;
            }
          }
        } else {
          if(!field.optional && !field.detail) {
            console.log(json);
            console.log(`Failed to instantiate ${(<any>this.constructor).name}: Non-optional, non-detail field ${fieldName} was missing`);
            throw new Error(`Failed to instantiate ${(<any>this.constructor).name}: Non-optional, non-detail field ${fieldName} was missing`);
          }
        }
      }
    }
  }

  /**
   * Serializes the model to JSON
   * @returns {any} The serialized JSON
   */
  toJSON(): any {
    let obj = {};
    for(let field of this.modelFields) {
      const fieldName = field.name;
      if(this.hasOwnProperty(field.name)) {
        switch(field.type) {
          case FieldType.PRIMITIVE:
            obj[fieldName] = this[fieldName];
            break;
          case FieldType.DATE:
            obj[fieldName] = this[fieldName].toString();
            break;
          case FieldType.MODEL:
            obj[fieldName] = this[fieldName].toJSON();
            break;
          case FieldType.MODEL_LIST:
            obj[fieldName] = this[fieldName].map(model => model.toJSON());
            break;
          case FieldType.MODEL_DICT: {
            let tmpObj = {};
            for(let key of this[fieldName]) {
              tmpObj[key] = this[fieldName][key].toJSON();
            }
            obj[fieldName] = tmpObj;
          }
        }
      } else {
        if(!field.optional && !field.detail) {
          console.log(`Failed to serialize ${(<any>this.constructor).name}: Non-optional, non-detail field ${fieldName} was missing`);
          throw new Error(`Failed to serialize ${(<any>this.constructor).name}: Non-optional, non-detail field ${fieldName} was missing`);

        }
      }
    }
    return obj;
  }
}