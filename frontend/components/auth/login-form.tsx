"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { Controller, useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod/v4-mini"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card"
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { PasswordInput } from "./password-input"

const formSchema = z.object({
  email: z.email(),
  password: z.string()
})

export function LoginForm() {

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: ""
    }
  })

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const url = "http://localhost:8080/api/v1/auth/login"

    try {
      const resp = await fetch(url, {
        method: 'POST',
        credentials: "include",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({email: data.email, password: data.password})
      })

      toast.success("Login successfull!")
    } catch (err) {
      toast.error("Login failed")
      console.log(err)
    }
  }

    return (
      <Card className="w-full sm:max-w-md">
        <CardHeader>
          <CardTitle>Login Page</CardTitle>
        </CardHeader>
        <CardContent>
          <form id="form-rhf-login" onSubmit={form.handleSubmit(onSubmit)}>
            <FieldGroup>
              <Controller
                name="email"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="form-rhf-demo-title">
                      Email
                    </FieldLabel>
                    <Input
                      {...field}
                      id="form-rhf-title"
                      aria-invalid={fieldState.invalid}
                      placeholder="Enter your email address."
                      autoComplete="off"
                    />
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
            />
            <Controller
              name="password"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor="form-rhf-demo-title">
                   Password
                  </FieldLabel>
                  <PasswordInput
                    {...field}
                    id="form-rhf-title"
                    aria-invalid={fieldState.invalid}
                    placeholder="Enter your password"
                    autoComplete="off"
                  />
                  {fieldState.invalid && (
                    <FieldError errors={[fieldState.error]} />
                  )}
                </Field>
              )}
            />
            </FieldGroup>
          </form>
        </CardContent>
        <CardFooter>
          <Field orientation="horizontal">
            <Button type="button" variant="outline" onClick={() => form.reset()}>
              Reset
            </Button>
            <Button type="submit" form="form-rhf-login">
              Submit
            </Button>
          </Field>
        </CardFooter>
      </Card>
    )}
