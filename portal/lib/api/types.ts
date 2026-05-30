export type LoginResponse = {
  access_token: string;
  token_type?: string;
  expires_in?: number;
  success?: boolean;
  data?: UserProfile;
};

export type UserProfile = {
  id: number;
  username: string;
  display_name: string;
  email: string;
  role: number;
  group: string;
  status: number;
  quota?: number;
  used_quota?: number;
  invite_code?: string;
};

export type UserSelfResponse = {
  success: boolean;
  data: UserProfile;
};

export type RegisterRequest = {
  email: string;
  username: string;
  password: string;
  invite_code?: string;
};
