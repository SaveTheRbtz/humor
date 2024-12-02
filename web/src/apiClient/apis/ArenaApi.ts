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


import * as runtime from '../runtime';
import type {
  ArenaRateChoicesBody,
  RpcStatus,
  V1GetChoicesResponse,
  V1GetLeaderboardResponse,
  V1GetTopJokesResponse,
} from '../models/index';
import {
    ArenaRateChoicesBodyFromJSON,
    ArenaRateChoicesBodyToJSON,
    RpcStatusFromJSON,
    RpcStatusToJSON,
    V1GetChoicesResponseFromJSON,
    V1GetChoicesResponseToJSON,
    V1GetLeaderboardResponseFromJSON,
    V1GetLeaderboardResponseToJSON,
    V1GetTopJokesResponseFromJSON,
    V1GetTopJokesResponseToJSON,
} from '../models/index';

export interface ArenaGetChoicesRequest {
    sessionId?: string;
}

export interface ArenaRateChoicesRequest {
    id: string;
    body: ArenaRateChoicesBody;
}

/**
 * 
 */
export class ArenaApi extends runtime.BaseAPI {

    /**
     * Retrieves a pair of jokes for comparison.
     */
    async arenaGetChoicesRaw(requestParameters: ArenaGetChoicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<V1GetChoicesResponse>> {
        const queryParameters: any = {};

        if (requestParameters['sessionId'] != null) {
            queryParameters['sessionId'] = requestParameters['sessionId'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/choice`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => V1GetChoicesResponseFromJSON(jsonValue));
    }

    /**
     * Retrieves a pair of jokes for comparison.
     */
    async arenaGetChoices(requestParameters: ArenaGetChoicesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<V1GetChoicesResponse> {
        const response = await this.arenaGetChoicesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Gets the leaderboard of joke models.
     */
    async arenaGetLeaderboardRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<V1GetLeaderboardResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/leaderboard`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => V1GetLeaderboardResponseFromJSON(jsonValue));
    }

    /**
     * Gets the leaderboard of joke models.
     */
    async arenaGetLeaderboard(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<V1GetLeaderboardResponse> {
        const response = await this.arenaGetLeaderboardRaw(initOverrides);
        return await response.value();
    }

    /**
     * Gets the top jokes.
     */
    async arenaGetTopJokesRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<V1GetTopJokesResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/top-jokes`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => V1GetTopJokesResponseFromJSON(jsonValue));
    }

    /**
     * Gets the top jokes.
     */
    async arenaGetTopJokes(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<V1GetTopJokesResponse> {
        const response = await this.arenaGetTopJokesRaw(initOverrides);
        return await response.value();
    }

    /**
     * Submits the user\'s choice between two jokes.
     */
    async arenaRateChoicesRaw(requestParameters: ArenaRateChoicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling arenaRateChoices().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling arenaRateChoices().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/choice/{id}/rate`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: ArenaRateChoicesBodyToJSON(requestParameters['body']),
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * Submits the user\'s choice between two jokes.
     */
    async arenaRateChoices(requestParameters: ArenaRateChoicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.arenaRateChoicesRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
