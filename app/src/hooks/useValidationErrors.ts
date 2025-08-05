import { useState } from 'react';

export interface ValidationErrors {
  [field: string]: string;
}

export function useValidationErrors() {
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>({});

  const setErrorsFromException = (error: any) => {
    // Réinitialiser les erreurs
    setValidationErrors({});
    
    // Extraire les erreurs de validation si elles existent
    if (error?.validationErrors && typeof error.validationErrors === 'object') {
      setValidationErrors(error.validationErrors);
      return true; // Indique qu'on a trouvé des erreurs de validation
    }
    
    return false; // Pas d'erreurs de validation spécifiques
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