/* tslint:disable */
/* eslint-disable */
/**
 * proto/server.proto
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: version not set
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
/**
 * LeaderboardEntry contains the model name and its Bradley-Terry rating.
 * @export
 * @interface V1LeaderboardEntry
 */
export interface V1LeaderboardEntry {
    /**
     * Public model name.
     * @type {string}
     * @memberof V1LeaderboardEntry
     */
    model?: string;
    /**
     * Bradley-Terry rating of the model.
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    bradleyterrScore?: number;
}

/**
 * Check if a given object implements the V1LeaderboardEntry interface.
 */
export function instanceOfV1LeaderboardEntry(value: object): value is V1LeaderboardEntry {
    return true;
}

export function V1LeaderboardEntryFromJSON(json: any): V1LeaderboardEntry {
    return V1LeaderboardEntryFromJSONTyped(json, false);
}

export function V1LeaderboardEntryFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1LeaderboardEntry {
    if (json == null) {
        return json;
    }
    return {
        
        'model': json['model'] == null ? undefined : json['model'],
        'bradleyterrScore': json['bradleyterrScore'] == null ? undefined : json['bradleyterrScore'],
    };
}

export function V1LeaderboardEntryToJSON(value?: V1LeaderboardEntry | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'model': value['model'],
        'bradleyterrScore': value['bradleyterrScore'],
    };
}

