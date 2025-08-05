import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useValidationErrors } from "@/hooks/use-validation-errors"
import { setupAPI } from "@/lib/api"
import { cn } from "@/lib/utils"
import { SetupAdminUserRequest } from "@/types/setup"
import { useForm } from '@tanstack/react-form'
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useRouter } from 'next/navigation'
import { useState } from 'react'

export function SetupForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [setupError, setSetupError] = useState<string | null>(null);
  const { validationErrors, setErrorsFromException, getFieldError, clearErrors } = useValidationErrors();

  const setupMutation = useMutation({
    mutationFn: setupAPI.createAdminUser,
    onSuccess: () => {
      // Refresh the health status so the app knows that the first-time setup is done
      queryClient.invalidateQueries({ queryKey: ['health'] });
      router.push('/login');
    },
    onError: (error: any) => {
      // Essayer d'extraire les erreurs de validation
      const hasValidationErrors = setErrorsFromException(error);
      
      // Si pas d'erreurs de validation spécifiques, afficher l'erreur générale
      if (!hasValidationErrors) {
        setSetupError(error.message || 'An error occurred while creating the account');
      }
    }
  });

  const form = useForm({
    defaultValues: {
      first_name: '',
      last_name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
    onSubmit: async ({ value }) => {
      setSetupError(null);
      clearErrors();
      
      if (value.password !== value.confirmPassword) {
        setSetupError('Passwords do not match');
        return;
      }

      const userData: SetupAdminUserRequest = {
        first_name: value.first_name,
        last_name: value.last_name || undefined,
        email: value.email,
        password: value.password,
      };

      try {
        await setupMutation.mutateAsync(userData);
      } catch (error) {
        console.error('Setup failed:', error);
      }
    },
  });

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="backdrop-blur-lg bg-card/80 dark:bg-card/60 border-2 border-border/50 shadow-2xl shadow-black/10 dark:shadow-black/30">
        <CardHeader className="text-center pb-8 pt-8">
          <CardTitle className="text-2xl font-jersey tracking-wide bg-gradient-to-r from-foreground via-foreground/90 to-foreground/70 bg-clip-text text-transparent">
            FIRST SETUP
          </CardTitle>
          <CardDescription className="text-muted-foreground/80 mt-2">
            Create your administrator account to get started
          </CardDescription>
          {/* Decorative line */}
          <div className="mx-auto w-16 h-0.5 bg-gradient-to-r from-transparent via-primary/50 to-transparent mt-4" />
        </CardHeader>
        <CardContent className="px-8 pb-8">
          <form onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              void form.handleSubmit();
            }}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <div className="grid gap-3">
                    <Label htmlFor="first_name" className="text-foreground/90 font-medium">
                      First Name
                    </Label>
                    <form.Field
                      name="first_name"
                      validators={{
                        onChange: ({ value }) => {
                          if (!value) return 'First name is required'
                          if (value.length < 2) return 'First name must be at least 2 characters'
                          return undefined
                        }
                      }}
                    >
                      {(field) => {
                        const serverError = getFieldError('first_name');
                        return (
                          <>
                            <Input
                              id={field.name}
                              type="text"
                              placeholder="John"
                              required
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) => field.handleChange(e.target.value)}
                              className="h-12 bg-background/50 border-border/60 focus:border-primary/60 focus:ring-primary/20 transition-all duration-300 hover:border-border/80"
                            />
                            {field.state.meta.isTouched && !field.state.meta.isValid && (
                              <p className="text-sm text-destructive mt-1">{field.state.meta.errors.join(', ')}</p>
                            )}
                            {serverError && (
                              <p className="text-sm text-destructive mt-1">{serverError}</p>
                            )}
                          </>
                        );
                      }}
                    </form.Field>
                  </div>

                  <div className="grid gap-3">
                    <Label htmlFor="last_name" className="text-foreground/90 font-medium">
                      Last Name <span className="text-muted-foreground">(optional)</span>
                    </Label>
                    <form.Field
                      name="last_name"
                    >
                      {(field) => {
                        const serverError = getFieldError('last_name');
                        return (
                          <>
                            <Input
                              id={field.name}
                              type="text"
                              placeholder="Doe"
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) => field.handleChange(e.target.value)}
                              className="h-12 bg-background/50 border-border/60 focus:border-primary/60 focus:ring-primary/20 transition-all duration-300 hover:border-border/80"
                            />
                            {serverError && (
                              <p className="text-sm text-destructive mt-1">{serverError}</p>
                            )}
                          </>
                        );
                      }}
                    </form.Field>
                  </div>

                  <div className="grid gap-3">
                    <Label htmlFor="email" className="text-foreground/90 font-medium">
                      Email
                    </Label>
                    <form.Field
                      name="email"
                      validators={{
                        onChange: ({ value }) => {
                          if (!value) return 'Email is required'
                          if (!/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i.test(value)) {
                            return 'Invalid email address'
                          }
                          return undefined
                        }
                      }}
                    >
                      {(field) => {
                        const serverError = getFieldError('email');
                        return (
                          <>
                            <Input
                              id={field.name}
                              type="email"
                              placeholder="admin@company.com"
                              required
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) => field.handleChange(e.target.value)}
                              className="h-12 bg-background/50 border-border/60 focus:border-primary/60 focus:ring-primary/20 transition-all duration-300 hover:border-border/80"
                            />
                            {field.state.meta.isTouched && !field.state.meta.isValid && (
                              <p className="text-sm text-destructive mt-1">{field.state.meta.errors.join(', ')}</p>
                            )}
                            {serverError && (
                              <p className="text-sm text-destructive mt-1">{serverError}</p>
                            )}
                          </>
                        );
                      }}
                    </form.Field>
                  </div>

                  <div className="grid gap-3">
                    <Label htmlFor="password" className="text-foreground/90 font-medium">
                      Password
                    </Label>
                    <form.Field
                      name="password"
                      validators={{
                        onChange: ({ value }) => {
                          if (!value) return 'Password is required'
                          if (value.length < 8) return 'Password must be at least 8 characters'
                          return undefined
                        }
                      }}
                    >
                      {(field) => {
                        const serverError = getFieldError('password');
                        return (
                          <>
                            <Input 
                              id={field.name}
                              type="password" 
                              required 
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) => field.handleChange(e.target.value)}
                              className="h-12 bg-background/50 border-border/60 focus:border-primary/60 focus:ring-primary/20 transition-all duration-300 hover:border-border/80"
                            />
                            {field.state.meta.isTouched && !field.state.meta.isValid && (
                              <p className="text-sm text-destructive mt-1">{field.state.meta.errors.join(', ')}</p>
                            )}
                            {serverError && (
                              <p className="text-sm text-destructive mt-1">{serverError}</p>
                            )}
                          </>
                        );
                      }}
                    </form.Field>
                  </div>

                  <div className="grid gap-3">
                    <Label htmlFor="confirmPassword" className="text-foreground/90 font-medium">
                      Confirm Password
                    </Label>
                    <form.Field
                      name="confirmPassword"
                      validators={{
                        onChange: ({ value }) => {
                          if (!value) return 'Please confirm your password'
                          return undefined
                        }
                      }}
                    >
                      {(field) => {
                        const serverError = getFieldError('confirmPassword');
                        return (
                          <>
                            <Input 
                              id={field.name}
                              type="password" 
                              required 
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) => field.handleChange(e.target.value)}
                              className="h-12 bg-background/50 border-border/60 focus:border-primary/60 focus:ring-primary/20 transition-all duration-300 hover:border-border/80"
                            />
                            {field.state.meta.isTouched && !field.state.meta.isValid && (
                              <p className="text-sm text-destructive mt-1">{field.state.meta.errors.join(', ')}</p>
                            )}
                            {serverError && (
                              <p className="text-sm text-destructive mt-1">{serverError}</p>
                            )}
                          </>
                        );
                      }}
                    </form.Field>
                  </div>

                  {setupError && (
                    <p className="text-sm text-destructive">{setupError}</p>
                  )}

                  <Button 
                    type="submit" 
                    disabled={setupMutation.isPending}
                    className="w-full h-12 font-jersey font-bold text-lg tracking-wide mt-2 bg-gradient-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]"
                  >
                    {setupMutation.isPending ? 'CREATING ACCOUNT...' : 'CREATE ADMIN ACCOUNT'}
                  </Button>
                </div>
              </div>
            </form>
        </CardContent>
      </Card>
    </div>
  )
}