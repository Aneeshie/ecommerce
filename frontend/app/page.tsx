"use client"
import { useUserStore } from "@/stores/auth-store";

export default function Home() {

  const user = useUserStore((state) => state.user)

  return (
    <div>
      {user ? (
        <>
          <p>{user.name}</p>
          <p>{user.email}</p>
          <p>{user.role}</p>
        </>
      ) : (
        <p>Not logged in</p>
      )}
    </div>
  );
}
