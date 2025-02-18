import { useRouter } from 'next/router';
import ChatRoom from '../components/chat-room';

const RoomPage = () => {
  const router = useRouter();
  const { room } = router.query;

  if (!room) {
    return <div>Loading...</div>;
  }

  return <ChatRoom room={room as string} />;
};

export default RoomPage;