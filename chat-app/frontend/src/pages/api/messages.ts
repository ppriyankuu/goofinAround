import type { NextApiRequest, NextApiResponse } from 'next';
import axios from 'axios';

interface Message {
  id: string;
  user: string;
  content: string;
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Message[]>) {
  const { room } = req.query;

  try {
    const response = await axios.get(`http://localhost:8080/rooms/${room}/messages`);
    res.status(200).json(response.data);
  } catch (error) {
    res.status(500).json([]);
  }
}