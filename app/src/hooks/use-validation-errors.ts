import type { ApiError } from '@/types/api';
import { useState } from 'react';

export interface ValidationErrors {
  [field: string]: string;
}

export function useValidationErrors() {
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>({});

  const setErrorsFromException = (error: ApiError) => {
    setValidationErrors({});
    if (error.validationErrors && typeof error.validationErrors === 'object') {
      setValidationErrors(error.validationErrors);
      return true;
    }
    return false;
  };

  const getFieldError = (fieldName: string): string | undefined => {
    return validationErrors[fieldName];
  };

  const clearErrors = () => {
    setValidationErrors({});
  };

  const hasErrors = () => {
    return Object.keys(validationErrors).length > 0;
  };

  return {
    validationErrors,
    setErrorsFromException,
    getFieldError,
    clearErrors,
    hasErrors,
  };
}