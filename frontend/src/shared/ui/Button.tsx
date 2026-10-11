import type { ButtonHTMLAttributes } from "react";
import { Spinner } from "./Spinner";
import styles from "./Button.module.css";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "danger" | "ghost";
  size?: "sm" | "md";
  loading?: boolean;
};

export function Button({
  variant = "secondary",
  size = "md",
  loading = false,
  disabled,
  className,
  children,
  type = "button",
  ...rest
}: Props) {
  const classes = [styles.button, styles[variant], styles[size], className].filter(Boolean).join(" ");
  return (
    <button type={type} className={classes} disabled={disabled || loading} {...rest}>
      {loading && <Spinner size={14} />}
      {children}
    </button>
  );
}
