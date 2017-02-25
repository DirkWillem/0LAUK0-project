import { Model, SubModelClass } from "./model";

/**
 * Represents an addition to a collection
 */
export interface CollectionAddition<M extends Model> {
  action: "added";
  addedEntity: M;
}

/**
 * Represents an update to a collection
 */
export interface CollectionUpdate<M extends Model> {
  action: "updated";
  updatedEntity: M;
}

/**
 * Represents a removal from a collection
 */
export interface CollectionRemoval {
  action: "deleted";
  id: number;
}

/**
 * Contains all possible collection mutation actions
 */
export type CollectionMutation<M extends Model> = CollectionAddition<M> | CollectionUpdate<M> | CollectionRemoval;

/**
 * Applies a collection mutation to a collection
 * @param collection - Collection to apply the mutation to
 * @param mutation - The mutation to apply
 * @returns {M[]} The updated collection
 */
export function applyUpdateToCollection<M extends Model>(collection: M[], mutation: CollectionMutation<M>): M[] {
  switch(mutation.action) {
    case "added":
      return [...collection, mutation.addedEntity];
    case "updated":
      return collection.map(entity => entity.id == mutation.updatedEntity.id ? mutation.updatedEntity : entity);
    case "deleted":
      return collection.filter(entity => entity.id != mutation.id);
  }
}