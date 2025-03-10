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
 * GetChoicesResponse is a response to the GetChoicesRequest.
 * @export
 * @interface V1GetChoicesResponse
 */
export interface V1GetChoicesResponse {
    /**
     * Unique identifier for the joke pair.
     * @type {string}
     * @memberof V1GetChoicesResponse
     */
    id?: string;
    /**
     * Theme for both jokes.
     * @type {string}
     * @memberof V1GetChoicesResponse
     */
    theme?: string;
    /**
     * Text of the left joke.
     * @type {string}
     * @memberof V1GetChoicesResponse
     */
    leftJoke?: string;
    /**
     * Text of the right joke.
     * @type {string}
     * @memberof V1GetChoicesResponse
     */
    rightJoke?: string;
}

/**
 * Check if a given object implements the V1GetChoicesResponse interface.
 */
export function instanceOfV1GetChoicesResponse(value: object): value is V1GetChoicesResponse {
    return true;
}

export function V1GetChoicesResponseFromJSON(json: any): V1GetChoicesResponse {
    return V1GetChoicesResponseFromJSONTyped(json, false);
}

export function V1GetChoicesResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1GetChoicesResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'] == null ? undefined : json['id'],
        'theme': json['theme'] == null ? undefined : json['theme'],
        'leftJoke': json['leftJoke'] == null ? undefined : json['leftJoke'],
        'rightJoke': json['rightJoke'] == null ? undefined : json['rightJoke'],
    };
}

export function V1GetChoicesResponseToJSON(value?: V1GetChoicesResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'theme': value['theme'],
        'leftJoke': value['leftJoke'],
        'rightJoke': value['rightJoke'],
    };
}

