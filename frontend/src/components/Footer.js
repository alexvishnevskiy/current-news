import * as React from 'react';
import Paper from '@mui/material/Paper';
import Stack from '@mui/material/Stack';
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import Link from '@mui/material/Link';
import { Typography } from '@mui/material';


export default function Footer() {
    return (
      <Paper sx={{marginTop: 'calc(3% + 10px)',
      width: '100%',
      }} component="footer" square variant="outlined">
        <Stack
        direction="row"
        justifyContent="center"
        // alignItems="center"
        spacing={2}
        >
            <Typography>
                © 2022 Alexander Vishnevskiy
            </Typography>
            <Link href="https://github.com/alexvishnevskiy" className={"github_link"} target="_blank" >
                <GitHubIcon/>
            </Link>
            <Link href="https://www.linkedin.com/in/alexandervishnevskiy/" className={"github_link"} target="_blank" >
                <LinkedInIcon/>
            </Link>
        </Stack>
      </Paper>
    );
  }