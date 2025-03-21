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
     * Total votes for the model.
     * @type {string}
     * @memberof V1LeaderboardEntry
     */
    votes?: string;
    /**
     * 
     * @type {string}
     * @memberof V1LeaderboardEntry
     */
    votesGood?: string;
    /**
     * 
     * @type {string}
     * @memberof V1LeaderboardEntry
     */
    votesBad?: string;
    /**
     * Newman Score of the model.
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    newmanScore?: number;
    /**
     * 
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    newmanCILower?: number;
    /**
     * 
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    newmanCIUpper?: number;
    /**
     * Elo Score of the model.
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    eloScore?: number;
    /**
     * 
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    eloCILower?: number;
    /**
     * 
     * @type {number}
     * @memberof V1LeaderboardEntry
     */
    eloCIUpper?: number;
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
        'votes': json['votes'] == null ? undefined : json['votes'],
        'votesGood': json['votesGood'] == null ? undefined : json['votesGood'],
        'votesBad': json['votesBad'] == null ? undefined : json['votesBad'],
        'newmanScore': json['newmanScore'] == null ? undefined : json['newmanScore'],
        'newmanCILower': json['newmanCILower'] == null ? undefined : json['newmanCILower'],
        'newmanCIUpper': json['newmanCIUpper'] == null ? undefined : json['newmanCIUpper'],
        'eloScore': json['eloScore'] == null ? undefined : json['eloScore'],
        'eloCILower': json['eloCILower'] == null ? undefined : json['eloCILower'],
        'eloCIUpper': json['eloCIUpper'] == null ? undefined : json['eloCIUpper'],
    };
}

export function V1LeaderboardEntryToJSON(value?: V1LeaderboardEntry | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'model': value['model'],
        'votes': value['votes'],
        'votesGood': value['votesGood'],
        'votesBad': value['votesBad'],
        'newmanScore': value['newmanScore'],
        'newmanCILower': value['newmanCILower'],
        'newmanCIUpper': value['newmanCIUpper'],
        'eloScore': value['eloScore'],
        'eloCILower': value['eloCILower'],
        'eloCIUpper': value['eloCIUpper'],
    };
}

