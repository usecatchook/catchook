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
import { useAuth } from "@/contexts/auth-context"
import { cn } from "@/lib/utils"
import { useForm } from '@tanstack/react-form'

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { login, isLoggingIn, loginError, getFieldError } = useAuth();

  const form = useForm({
    defaultValues: {
      email: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
      try {
        await login(value);
        // La redirection est gérée automatiquement par la fonction login
      } catch (error) {
        console.error('Login failed:', error);
      }
    },
  });

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="backdrop-blur-lg bg-card/80 dark:bg-card/60 border-2 border-border/50 shadow-2xl shadow-black/10 dark:shadow-black/30">
        <CardHeader className="text-center pb-8 pt-8">
          <CardTitle className="text-2xl font-jersey tracking-wide bg-gradient-to-r from-foreground via-foreground/90 to-foreground/70 bg-clip-text text-transparent">
            WELCOME TO CATCHOOK
          </CardTitle>
          <CardDescription className="text-muted-foreground/80 mt-2">
            Login with your email and password
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
                              placeholder="john@doe.com"
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
                    <div className="flex items-center">
                      <Label htmlFor="password" className="text-foreground/90 font-medium">
                        Password
                      </Label>
                      <a
                        href="#"
                        className="ml-auto text-sm text-muted-foreground hover:text-primary underline-offset-4 hover:underline transition-colors duration-200"
                      >
                        Forgot your password?
                      </a>
                    </div>
                    <form.Field
                      name="password"
                      validators={{
                        onChange: ({ value }) => {
                          if (!value) return 'Password is required'
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
                  {loginError && (
                    <p className="text-sm text-destructive">{loginError}</p>
                  )}
                  <Button 
                    type="submit" 
                    disabled={isLoggingIn}
                    className="w-full h-12 font-jersey font-bold text-lg tracking-wide mt-2 bg-gradient-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]"
                  >
                    {isLoggingIn ? 'LOGGING IN...' : 'LOGIN'}
                  </Button>
                </div>
              </div>
            </form>
        </CardContent>
      </Card>
    </div>
  )
}
