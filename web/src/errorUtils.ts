import { ResponseError, FetchError } from './apiClient';

export async function getErrorMessage(err: any): Promise<string> {
  console.error('Error details:', err);

  let errorMessage = '';

  if (err instanceof ResponseError) {
    const response = err.response;
    errorMessage = `HTTP ${response.status}: ${response.statusText}`;

    try {
      const contentType = response.headers.get('Content-Type');
      let errorBody: any;

      if (contentType && contentType.includes('application/json')) {
        errorBody = await response.json();
      } else {
        errorBody = await response.text();
      }

      if (errorBody) {
        if (typeof errorBody === 'string') {
          try {
            errorBody = JSON.parse(errorBody);
          } catch (parseErr) {
          }
        }

        if (errorBody.message) {
          errorMessage += ` - ${errorBody.message}`;
        } else {
          errorMessage += ` - ${JSON.stringify(errorBody)}`;
        }
      }
    } catch (parseErr) {
      console.error('Failed to parse error body:', parseErr);
    }
  } else if (err instanceof FetchError) {
    errorMessage = `Network error: ${err.message}`;
  } else if (err.message) {
    errorMessage = `Error: ${err.message}`;
  } else {
    errorMessage = `Error: ${String(err)}`;
  }

  return errorMessage;
}