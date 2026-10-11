export type PaymentSettings = {
  details: string;
  qr_url: string;
  currency: string;
  updated_at: string;
};

export type PurchaseStatus = "pending" | "approved" | "rejected";

export type ContactType = "phone" | "whatsapp" | "telegram";

export type AdminPurchase = {
  id: string;
  book_id: string;
  amount: number;
  currency: string;
  status: PurchaseStatus;
  created_at: string;
  reviewed_at?: string;
  reviewed_by?: string;
  book_title: string;
  buyer_name: string;
  buyer_email: string;
  contact_type: ContactType;
  contact_value: string;
  receipt_mime: string;
};
