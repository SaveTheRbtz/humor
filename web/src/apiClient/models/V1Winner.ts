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


/**
 * Enumeration of possible user choices.
 * 
 *  - UNSPECIFIED: Default unspecified value.
 *  - NONE: User didn't like either joke.
 *  - LEFT: User chose the left joke.
 *  - RIGHT: User chose the right joke.
 *  - BOTH: User liked both jokes equally.
 * @export
 */
export const V1Winner = {
    Unspecified: 'UNSPECIFIED',
    None: 'NONE',
    Left: 'LEFT',
    Right: 'RIGHT',
    Both: 'BOTH'
} as const;
export type V1Winner = typeof V1Winner[keyof typeof V1Winner];


export function instanceOfV1Winner(value: any): boolean {
    for (const key in V1Winner) {
        if (Object.prototype.hasOwnProperty.call(V1Winner, key)) {
            if (V1Winner[key as keyof typeof V1Winner] === value) {
                return true;
            }
        }
    }
    return false;
}

export function V1WinnerFromJSON(json: any): V1Winner {
    return V1WinnerFromJSONTyped(json, false);
}

export function V1WinnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1Winner {
    return json as V1Winner;
}

export function V1WinnerToJSON(value?: V1Winner | null): any {
    return value as any;
}

