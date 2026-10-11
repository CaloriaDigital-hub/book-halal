export type Requisites = {
  details: string;
  qr_url: string;
  currency: string;
};

export type PurchaseStatus = "pending" | "approved" | "rejected";

export type MyPurchase = {
  id: string;
  amount: number;
  currency: string;
  status: PurchaseStatus;
  created_at: string;
  reviewed_at?: string;
};

export type ContactType = "whatsapp" | "telegram" | "phone";
