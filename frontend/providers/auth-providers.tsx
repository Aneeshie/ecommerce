"use client";

import { ReactNode, useEffect } from "react";
import { useUserStore } from "@/stores/auth-store";

interface Props {
  children: ReactNode;
}

const AuthProvider = ({ children }: Props) => {
  const {
    loading,
    setUser,
    clearUser,
    setLoading,
  } = useUserStore();

  useEffect(() => {
    const fetchUser = async () => {
      setLoading(true);

      try {
        const resp = await fetch("http://localhost:8080/api/v1/auth/me", {
          credentials: "include",
        });

        if (!resp.ok) {
          clearUser();
          return;
        }

        const user = await resp.json();
        setUser(user);
      } catch (err) {
        console.error(err);
        clearUser();
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [setUser, clearUser, setLoading]);

  if (loading) {
    return <div>Loading...</div>;
  }

  return <>{children}</>;
};

export default AuthProvider;
